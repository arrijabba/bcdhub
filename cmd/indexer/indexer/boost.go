package indexer

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/baking-bad/bcdhub/internal/config"
	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/contractparser/kinds"
	"github.com/baking-bad/bcdhub/internal/contractparser/meta"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/helpers"
	"github.com/baking-bad/bcdhub/internal/index"
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/mq"
	"github.com/baking-bad/bcdhub/internal/noderpc"
	"github.com/baking-bad/bcdhub/internal/parsers"
	"github.com/baking-bad/bcdhub/internal/parsers/operations"
	"github.com/baking-bad/bcdhub/internal/rollback"
	"github.com/pkg/errors"
)

var errBcdQuit = errors.New("bcd-quit")
var errRollback = errors.New("rollback")
var errSameLevel = errors.New("Same level")

// BoostIndexer -
type BoostIndexer struct {
	rpc             noderpc.INode
	es              elastic.IElastic
	externalIndexer index.Indexer
	interfaces      map[string]kinds.ContractKind
	messageQueue    mq.Mediator
	state           models.Block
	currentProtocol models.Protocol
	cfg             config.Config

	updateTicker        *time.Ticker
	stop                chan struct{}
	Network             string
	boost               bool
	skipDelegatorBlocks bool
	stopped             bool
}

func (bi *BoostIndexer) fetchExternalProtocols() error {
	logger.WithNetwork(bi.Network).Info("Fetching external protocols")
	var existingProtocols []models.Protocol
	if err := bi.es.GetByNetworkWithSort(bi.Network, "start_level", "desc", &existingProtocols); err != nil {
		return err
	}

	exists := make(map[string]bool)
	for _, existingProtocol := range existingProtocols {
		exists[existingProtocol.Hash] = true
	}

	extProtocols, err := bi.externalIndexer.GetProtocols()
	if err != nil {
		return err
	}

	protocols := make([]elastic.Model, 0)
	for i := range extProtocols {
		if _, ok := exists[extProtocols[i].Hash]; ok {
			continue
		}
		symLink, err := meta.GetProtoSymLink(extProtocols[i].Hash)
		if err != nil {
			return err
		}
		alias := extProtocols[i].Alias
		if alias == "" {
			alias = extProtocols[i].Hash[:8]
		}

		newProtocol := &models.Protocol{
			ID:         helpers.GenerateID(),
			Hash:       extProtocols[i].Hash,
			Alias:      alias,
			StartLevel: extProtocols[i].StartLevel,
			EndLevel:   extProtocols[i].LastLevel,
			SymLink:    symLink,
			Network:    bi.Network,
		}

		protocolConstants := models.Constants{}
		if newProtocol.StartLevel != newProtocol.EndLevel || newProtocol.EndLevel != 0 {
			constants, err := bi.rpc.GetNetworkConstants(extProtocols[i].StartLevel)
			if err != nil {
				return err
			}
			protocolConstants.CostPerByte = constants.CostPerByte
			protocolConstants.HardGasLimitPerOperation = constants.HardGasLimitPerOperation
			protocolConstants.HardStorageLimitPerOperation = constants.HardStorageLimitPerOperation
			protocolConstants.TimeBetweenBlocks = constants.TimeBetweenBlocks[0]
		}
		newProtocol.Constants = protocolConstants

		protocols = append(protocols, newProtocol)
		logger.WithNetwork(bi.Network).Infof("Fetched %s", alias)
	}

	return bi.es.BulkInsert(protocols)
}

// NewBoostIndexer -
func NewBoostIndexer(cfg config.Config, network string, opts ...BoostIndexerOption) (*BoostIndexer, error) {
	logger.WithNetwork(network).Info("Creating indexer object...")
	es := elastic.WaitNew(cfg.Elastic.URI, cfg.Elastic.Timeout)
	rpcProvider, ok := cfg.RPC[network]
	if !ok {
		return nil, errors.Errorf("Unknown network %s", network)
	}
	rpc := noderpc.NewWaitNodeRPC(
		rpcProvider.URI,
		noderpc.WithTimeout(time.Duration(rpcProvider.Timeout)*time.Second),
	)

	messageQueue := mq.New(cfg.RabbitMQ.URI, cfg.Indexer.ProjectName, cfg.Indexer.MQ.NeedPublisher, 10)

	interfaces, err := kinds.Load()
	if err != nil {
		return nil, err
	}

	bi := &BoostIndexer{
		Network:      network,
		rpc:          rpc,
		es:           es,
		messageQueue: messageQueue,
		stop:         make(chan struct{}),
		interfaces:   interfaces,
		cfg:          cfg,
	}

	for _, opt := range opts {
		opt(bi)
	}

	err = bi.init()
	return bi, err
}

func (bi *BoostIndexer) init() error {
	if err := bi.es.CreateIndexes(); err != nil {
		return err
	}

	if bi.boost {
		if err := bi.fetchExternalProtocols(); err != nil {
			return err
		}
	}

	currentState, err := bi.es.GetLastBlock(bi.Network)
	if err != nil {
		return err
	}
	bi.state = currentState
	logger.WithNetwork(bi.Network).Infof("Current indexer state: %d", currentState.Level)

	currentProtocol, err := bi.es.GetProtocol(bi.Network, "", currentState.Level)
	if err != nil {
		header, err := bi.rpc.GetHeader(helpers.MaxInt64(1, currentState.Level))
		if err != nil {
			return err
		}
		currentProtocol, err = createProtocol(bi.Network, header.Protocol, 0)
		if err != nil {
			return err
		}
	}
	bi.currentProtocol = currentProtocol
	logger.WithNetwork(bi.Network).Infof("Current network protocol: %s", currentProtocol.Hash)
	return nil
}

// Sync -
func (bi *BoostIndexer) Sync(wg *sync.WaitGroup) {
	defer wg.Done()

	bi.stopped = false
	localSentry := helpers.GetLocalSentry()
	helpers.SetLocalTagSentry(localSentry, "network", bi.Network)

	// First tick
	if err := bi.process(); err != nil {
		logger.Error(err)
		helpers.CatchErrorSentry(err)
	}
	if bi.stopped {
		return
	}

	everySecond := false
	duration := time.Duration(bi.currentProtocol.Constants.TimeBetweenBlocks) * time.Second
	if duration.Microseconds() <= 0 {
		duration = 10 * time.Second
	}
	bi.setUpdateTicker(0)
	for {
		select {
		case <-bi.stop:
			bi.stopped = true
			bi.messageQueue.Close()
			return
		case <-bi.updateTicker.C:
			if err := bi.process(); err != nil {
				if errors.Is(err, errSameLevel) {
					if !everySecond {
						everySecond = true
						bi.setUpdateTicker(5)
					}
					continue
				}
				logger.Error(err)
				helpers.CatchErrorSentry(err)
			}

			if everySecond {
				everySecond = false
				bi.setUpdateTicker(0)
			}
		}
	}
}

func (bi *BoostIndexer) setUpdateTicker(seconds int) {
	if bi.updateTicker != nil {
		bi.updateTicker.Stop()
	}
	var duration time.Duration
	if seconds == 0 {
		duration = time.Duration(bi.currentProtocol.Constants.TimeBetweenBlocks) * time.Second
		if duration.Microseconds() <= 0 {
			duration = 10 * time.Second
		}
	} else {
		duration = time.Duration(seconds) * time.Second
	}
	logger.WithNetwork(bi.Network).Infof("Data will be updated every %.0f seconds", duration.Seconds())
	bi.updateTicker = time.NewTicker(duration)
}

// Stop -
func (bi *BoostIndexer) Stop() {
	bi.stop <- struct{}{}
}

// Index -
func (bi *BoostIndexer) Index(levels []int64) error {
	if len(levels) == 0 {
		return nil
	}
	for _, level := range levels {
		helpers.SetTagSentry("level", fmt.Sprintf("%d", level))

		select {
		case <-bi.stop:
			bi.stopped = true
			bi.messageQueue.Close()
			return errBcdQuit
		default:
		}

		currentHead, err := bi.rpc.GetHeader(level)
		if err != nil {
			return err
		}

		if bi.state.Level > 0 && currentHead.Predecessor != bi.state.Hash && !bi.boost {
			return errRollback
		}

		logger.WithNetwork(bi.Network).Infof("indexing %d block", level)

		if currentHead.Protocol != bi.currentProtocol.Hash {
			logger.WithNetwork(bi.Network).Infof("New protocol detected: %s -> %s", bi.currentProtocol.Hash, currentHead.Protocol)
			migrationModels, err := bi.migrate(currentHead)
			if err != nil {
				return err
			}
			if err := bi.saveModels(migrationModels); err != nil {
				return err
			}
		}

		parsedModels, err := bi.getDataFromBlock(bi.Network, currentHead)
		if err != nil {
			return err
		}
		parsedModels = append(parsedModels, bi.createBlock(currentHead))

		if err := bi.saveModels(parsedModels); err != nil {
			return err
		}
	}
	return nil
}

// Rollback -
func (bi *BoostIndexer) Rollback() error {
	logger.WithNetwork(bi.Network).Warningf("Rollback from %d", bi.state.Level)

	lastLevel, err := bi.getLastRollbackBlock()
	if err != nil {
		return err
	}

	manager := rollback.NewManager(bi.es, bi.messageQueue, bi.rpc, bi.cfg.SharePath)
	if err := manager.Rollback(bi.state, lastLevel); err != nil {
		return err
	}

	helpers.CatchErrorSentry(errors.Errorf("[%s] Rollback from %d to %d", bi.Network, bi.state.Level, lastLevel))

	newState, err := bi.es.GetLastBlock(bi.Network)
	if err != nil {
		return err
	}
	bi.state = newState
	logger.WithNetwork(bi.Network).Infof("New indexer state: %d", bi.state.Level)
	logger.WithNetwork(bi.Network).Info("Rollback finished")
	return nil
}

func (bi *BoostIndexer) getLastRollbackBlock() (int64, error) {
	var lastLevel int64
	level := bi.state.Level

	for end := false; !end; level-- {
		headAtLevel, err := bi.rpc.GetHeader(level)
		if err != nil {
			return 0, err
		}

		block, err := bi.es.GetBlock(bi.Network, level)
		if err != nil {
			return 0, err
		}

		if block.Predecessor == headAtLevel.Predecessor {
			logger.WithNetwork(bi.Network).Warnf("Found equal predecessors at level: %d", block.Level)
			end = true
			lastLevel = block.Level - 1
		}
	}
	return lastLevel, nil
}

func (bi *BoostIndexer) getBoostBlocks(head noderpc.Header) ([]int64, error) {
	levels, err := bi.externalIndexer.GetContractOperationBlocks(bi.state.Level, head.Level, bi.skipDelegatorBlocks)
	if err != nil {
		return nil, err
	}

	protocols, err := bi.externalIndexer.GetProtocols()
	if err != nil {
		return nil, err
	}

	protocolLevels := make([]int64, 0)
	for i := range protocols {
		if protocols[i].StartLevel > bi.state.Level && protocols[i].StartLevel > 0 {
			protocolLevels = append(protocolLevels, protocols[i].StartLevel)
		}
	}

	result := helpers.Merge2ArraysInt64(levels, protocolLevels)
	return result, err
}

func (bi *BoostIndexer) process() error {
	head, err := bi.rpc.GetHead()
	if err != nil {
		return err
	}

	if !bi.state.ValidateChainID(head.ChainID) {
		return errors.Errorf("Invalid chain_id: %s (state) != %s (head)", bi.state.ChainID, head.ChainID)
	}

	logger.WithNetwork(bi.Network).Infof("Current node state: %d", head.Level)
	logger.WithNetwork(bi.Network).Infof("Current indexer state: %d", bi.state.Level)

	if head.Level > bi.state.Level {
		levels := make([]int64, 0)
		if bi.boost {
			levels, err = bi.getBoostBlocks(head)
			if err != nil {
				return err
			}
		} else {
			for i := bi.state.Level + 1; i <= head.Level; i++ {
				levels = append(levels, i)
			}
		}

		logger.WithNetwork(bi.Network).Infof("Found %d new levels", len(levels))

		if err := bi.Index(levels); err != nil {
			if errors.Is(err, errBcdQuit) {
				return nil
			}
			if errors.Is(err, errRollback) {
				if !time.Now().Add(time.Duration(-5) * time.Minute).After(head.Timestamp) { // Check that node is out of sync
					if err := bi.Rollback(); err != nil {
						return err
					}
				}
				return nil
			}
			return err
		}

		if bi.boost {
			bi.boost = false
		}
		logger.WithNetwork(bi.Network).Info("Synced")
		return nil
	} else if head.Level < bi.state.Level {
		if err := bi.Rollback(); err != nil {
			return err
		}
	}

	return errSameLevel
}

func (bi *BoostIndexer) createBlock(head noderpc.Header) *models.Block {
	newBlock := models.Block{
		ID:          helpers.GenerateID(),
		Network:     bi.Network,
		Hash:        head.Hash,
		Predecessor: head.Predecessor,
		Protocol:    head.Protocol,
		ChainID:     head.ChainID,
		Level:       head.Level,
		Timestamp:   head.Timestamp,
	}

	bi.state = newBlock
	return &newBlock
}

func (bi *BoostIndexer) saveModels(items []elastic.Model) error {
	logger.WithNetwork(bi.Network).Debugf("Found %d new models", len(items))
	if err := bi.es.BulkInsert(items); err != nil {
		return err
	}

	for i := range items {
		if err := bi.messageQueue.Send(items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (bi *BoostIndexer) getDataFromBlock(network string, head noderpc.Header) ([]elastic.Model, error) {
	if head.Level <= 1 {
		return nil, nil
	}
	data, err := bi.rpc.GetOperations(head.Level)
	if err != nil {
		return nil, err
	}

	parsedModels := make([]elastic.Model, 0)
	for _, opg := range data.Array() {
		parser := operations.NewGroup(operations.NewParseParams(
			bi.rpc,
			bi.es,
			operations.WithConstants(bi.currentProtocol.Constants),
			operations.WithHead(head),
			operations.WithIPFSGateways(bi.cfg.IPFSGateways),
			operations.WithInterfaces(bi.interfaces),
			operations.WithShareDirectory(bi.cfg.SharePath),
			operations.WithNetwork(network),
		))
		parsed, err := parser.Parse(opg)
		if err != nil {
			return nil, err
		}
		parsedModels = append(parsedModels, parsed...)
	}

	return parsedModels, nil
}

func (bi *BoostIndexer) migrate(head noderpc.Header) ([]elastic.Model, error) {
	updates := make([]elastic.Model, 0)
	newModels := make([]elastic.Model, 0)

	if bi.currentProtocol.EndLevel == 0 && head.Level > 1 {
		logger.WithNetwork(bi.Network).Infof("Finalizing the previous protocol: %s", bi.currentProtocol.Alias)
		bi.currentProtocol.EndLevel = head.Level - 1
		updates = append(updates, &bi.currentProtocol)
	}

	newProtocol, err := bi.es.GetProtocol(bi.Network, head.Protocol, head.Level)
	if err != nil {
		logger.Warning("%s", err)
		newProtocol, err = createProtocol(bi.Network, head.Protocol, head.Level)
		if err != nil {
			return nil, err
		}
	}

	if bi.Network == consts.Mainnet && head.Level == 1 {
		vestingMigrations, err := bi.vestingMigration(head)
		if err != nil {
			return nil, err
		}
		newModels = append(newModels, vestingMigrations...)
	} else {
		if bi.currentProtocol.SymLink == "" {
			return nil, errors.Errorf("[%s] Protocol should be initialized", bi.Network)
		}
		if newProtocol.SymLink != bi.currentProtocol.SymLink {
			migrations, migrationUpdates, err := bi.standartMigration(newProtocol, head)
			if err != nil {
				return nil, err
			}
			newModels = append(newModels, migrations...)
			if len(migrationUpdates) > 0 {
				updates = append(updates, migrationUpdates...)
			}
		} else {
			logger.WithNetwork(bi.Network).Infof("Same symlink %s for %s / %s",
				newProtocol.SymLink, bi.currentProtocol.Alias, newProtocol.Alias)
		}
	}

	bi.currentProtocol = newProtocol
	newModels = append(newModels, &newProtocol)

	if err := bi.es.BulkUpdate(updates); err != nil {
		return nil, err
	}

	bi.setUpdateTicker(0)
	logger.WithNetwork(bi.Network).Infof("Migration to %s is completed", bi.currentProtocol.Alias)
	return newModels, nil
}

func createProtocol(network, hash string, level int64) (protocol models.Protocol, err error) {
	logger.WithNetwork(network).Infof("Creating new protocol %s starting at %d", hash, level)
	protocol.SymLink, err = meta.GetProtoSymLink(hash)
	if err != nil {
		return
	}

	protocol.Alias = hash[:8]
	protocol.Network = network
	protocol.Hash = hash
	protocol.StartLevel = level
	protocol.ID = helpers.GenerateID()
	return
}

func (bi *BoostIndexer) standartMigration(newProtocol models.Protocol, head noderpc.Header) ([]elastic.Model, []elastic.Model, error) {
	logger.WithNetwork(bi.Network).Info("Try to find migrations...")
	contracts, err := bi.es.GetContracts(map[string]interface{}{
		"network": bi.Network,
	})
	if err != nil {
		return nil, nil, err
	}
	logger.WithNetwork(bi.Network).Infof("Now %d contracts are indexed", len(contracts))

	p := parsers.NewMigrationParser(bi.es, bi.cfg.SharePath)
	newModels := make([]elastic.Model, 0)
	newUpdates := make([]elastic.Model, 0)
	for i := range contracts {
		logger.WithNetwork(bi.Network).Infof("Migrate %s...", contracts[i].Address)
		script, err := bi.rpc.GetScriptJSON(contracts[i].Address, newProtocol.StartLevel)
		if err != nil {
			return nil, nil, err
		}

		createdModels, updates, err := p.Parse(script, contracts[i], bi.currentProtocol, newProtocol, head.Timestamp)
		if err != nil {
			return nil, nil, err
		}

		if len(createdModels) > 0 {
			newModels = append(newModels, createdModels...)
		}
		if len(updates) > 0 {
			newUpdates = append(newUpdates, updates...)
		}
	}
	return newModels, newUpdates, nil
}

func (bi *BoostIndexer) vestingMigration(head noderpc.Header) ([]elastic.Model, error) {
	addresses, err := bi.rpc.GetContractsByBlock(head.Level)
	if err != nil {
		return nil, err
	}

	p := parsers.NewVestingParser(bi.cfg.SharePath, bi.interfaces)

	parsedModels := make([]elastic.Model, 0)
	for _, address := range addresses {
		if !strings.HasPrefix(address, "KT") {
			continue
		}

		data, err := bi.rpc.GetContractJSON(address, head.Level)
		if err != nil {
			return nil, err
		}

		parsed, err := p.Parse(data, head, bi.Network, address)
		if err != nil {
			return nil, err
		}
		parsedModels = append(parsedModels, parsed...)
	}

	return parsedModels, nil
}
