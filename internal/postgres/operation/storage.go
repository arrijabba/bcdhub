package operation

import (
	"fmt"
	"time"

	"github.com/baking-bad/bcdhub/internal/bcd"
	constants "github.com/baking-bad/bcdhub/internal/bcd/consts"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/operation"
	"github.com/baking-bad/bcdhub/internal/postgres/core"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Storage -
type Storage struct {
	*core.Postgres
}

// NewStorage -
func NewStorage(es *core.Postgres) *Storage {
	return &Storage{es}
}

// GetOne -
func (storage *Storage) GetOne(hash string, counter int64, nonce *int64) (op operation.Operation, err error) {
	query := storage.DB.Table(models.DocOperations).
		Where("hash = ?", hash).
		Where("counter = ?", counter)

	if nonce == nil {
		query.Where("nonce IS NULL")
	} else {
		query.Where("nonce = ?", nonce)
	}

	err = query.First(&op).Error
	return
}

type opgForContract struct {
	Counter int64
	Hash    string
}

func (storage *Storage) getContractOPG(address, network string, size uint64, filters map[string]interface{}) (response []opgForContract, err error) {
	query := storage.DB.Table(models.DocOperations).Select("hash", "counter").
		Where("network = ?", network).
		Where(
			storage.DB.Where("source = ?", address).Or("destination = ?", address),
		)

	if err := prepareOperationFilters(query, filters); err != nil {
		return nil, err
	}

	limit := storage.GetPageSize(int64(size))
	query.Group("hash, counter, level").Order("level DESC").Limit(limit)

	err = query.Find(&response).Error
	return
}

func prepareOperationFilters(query *gorm.DB, filters map[string]interface{}) error {
	for k, v := range filters {
		if v != "" {
			switch k {
			case "from":
				query.Where("timestamp >= ?", v)
			case "to":
				query.Where("timestamp <= ?", v)
			case "entrypoints":
				query.Where("entrypoint IN ?", v)
			case "last_id":
				query.Where("id < ?", v)
			case "status":
				query.Where("status IN ?", v)
			default:
				return errors.Errorf("Unknown operation filter: %s %v", k, v)
			}
		}
	}
	return nil
}

// GetByContract -
func (storage *Storage) GetByContract(network, address string, size uint64, filters map[string]interface{}) (po operation.Pageable, err error) {
	opg, err := storage.getContractOPG(address, network, size, filters)
	if err != nil {
		return
	}
	if len(opg) == 0 {
		return
	}

	query := storage.DB.Table(models.DocOperations).Where("network = ?", network)

	subQuery := storage.DB.Where(
		storage.DB.Where("hash = ?", opg[0].Hash).Where("counter = ?", opg[0].Counter),
	)
	for i := 1; i < len(opg); i++ {
		subQuery.Or(
			storage.DB.Where("hash = ?", opg[i].Hash).Where("counter = ?", opg[i].Counter),
		)
	}
	query.Where(subQuery)

	addOperationSorting(query)

	if err = query.Find(&po.Operations).Error; err != nil {
		return
	}

	if len(po.Operations) == 0 {
		return
	}

	lastID := po.Operations[0].ID
	for _, op := range po.Operations[1:] {
		if op.ID > lastID {
			continue
		}
		lastID = op.ID
	}
	po.LastID = fmt.Sprintf("%d", lastID)
	return
}

// Last -
func (storage *Storage) Last(network, address string, indexedTime int64) (op operation.Operation, err error) {
	err = storage.DB.Table(models.DocOperations).
		Where("network = ?", network).
		Where("destination = ?", address).
		Where("id < ?", indexedTime).
		Where("status = ?", constants.Applied).
		Where("deffated_storage != ''").
		Order("id desc").
		First(&op).
		Error
	return
}

// Get -
func (storage *Storage) Get(filters map[string]interface{}, size int64, sort bool) (operations []operation.Operation, err error) {
	query := storage.DB.Table(models.DocOperations).Where(filters)

	if sort {
		addOperationSorting(query)
	}

	if size > 0 {
		query.Limit(storage.GetPageSize(size))
	}

	err = query.Find(&operations).Error
	return operations, err
}

// GetStats -
func (storage *Storage) GetStats(network, address string) (stats operation.Stats, err error) {
	query := storage.DB.Table(models.DocOperations).
		Select("MAX(timestamp) AS last_action, COUNT(*) as count").
		Where("network = ?", network).
		Where(
			storage.DB.Where("source = ?", address).Or("destination = ?", address),
		)

	err = query.Scan(&stats).Error
	return
}

// GetContract24HoursVolume -
func (storage *Storage) GetContract24HoursVolume(network, address string, entrypoints []string) (float64, error) {
	aDayAgo := time.Now().UTC().AddDate(0, 0, -1)

	var volume float64
	query := storage.DB.Table(models.DocOperations).
		Select("COALESCE(SUM(amount), 0)").
		Where("destination = ?", address).
		Where("network = ?", network).
		Where("status = ?", constants.Applied).
		Where("timestamp > ?", aDayAgo)

	if len(entrypoints) > 0 {
		query.Where("entrypoint IN ?", entrypoints)
	}

	err := query.Scan(&volume).Error
	return volume, err
}

type tokenStats struct {
	Destination string
	Entrypoint  string
	Gas         int64
	Count       int64
}

// GetTokensStats -
func (storage *Storage) GetTokensStats(network string, addresses, entrypoints []string) (map[string]operation.TokenUsageStats, error) {
	var stats []tokenStats
	query := storage.DB.Table(models.DocOperations).
		Select("destination, entrypoint, COUNT(*) as count, SUM(consumed_gas) AS gas").
		Where("network = ?", network)

	if len(addresses) > 0 {
		subQuery := storage.DB.Where("destination = ?", addresses[0])
		for i := 1; i < len(addresses); i++ {
			subQuery.Or("destination = ?", addresses[i])
		}
		query.Where(subQuery)
	}

	if len(entrypoints) > 0 {
		subQuery := storage.DB.Where("entrypoint = ?", entrypoints[0])
		for i := 1; i < len(entrypoints); i++ {
			subQuery.Or("entrypoint = ?", entrypoints[i])
		}
		query.Where(subQuery)
	}

	query.Group("destination, entrypoint")

	if err := query.Find(&stats).Error; err != nil {
		return nil, err
	}

	usageStats := make(map[string]operation.TokenUsageStats)
	for i := range stats {
		usage := operation.TokenMethodUsageStats{
			Count:       stats[i].Count,
			ConsumedGas: stats[i].Gas,
		}
		if _, ok := usageStats[stats[i].Destination]; !ok {
			usageStats[stats[i].Destination] = make(operation.TokenUsageStats)
		}
		usageStats[stats[i].Destination][stats[i].Entrypoint] = usage
	}

	return usageStats, nil
}

type operationAddresses struct {
	Source      string
	Destination string
}

// GetParticipatingContracts -
func (storage *Storage) GetParticipatingContracts(network string, fromLevel, toLevel int64) ([]string, error) {
	query := storage.DB.Table(models.DocOperations).
		Select("source, destination").
		Where("network = ?", network).
		Where("level <= ?", fromLevel).
		Where("level > ?", toLevel)

	var response []operationAddresses
	if err := query.Find(&response).Error; err != nil {
		return nil, err
	}

	exists := make(map[string]struct{})
	addresses := make([]string, 0)
	for _, op := range response {
		if _, ok := exists[op.Source]; !ok && bcd.IsContract(op.Source) {
			addresses = append(addresses, op.Source)
			exists[op.Source] = struct{}{}
		}
		if _, ok := exists[op.Destination]; !ok && bcd.IsContract(op.Destination) {
			addresses = append(addresses, op.Destination)
			exists[op.Destination] = struct{}{}
		}
	}

	return addresses, nil
}

// GetByIDs -
func (storage *Storage) GetByIDs(ids ...int64) (result []operation.Operation, err error) {
	err = storage.DB.Table(models.DocOperations).Order("id asc").Find(&result, ids).Error
	return
}

// GetDAppStats -
func (storage *Storage) GetDAppStats(network string, addresses []string, period string) (stats operation.DAppStats, err error) {
	query, err := getDAppQuery(storage.DB, network, addresses, period)
	if err != nil {
		return
	}

	if err = query.Select("COUNT(*) as calls, SUM(amount) as volume").First(&stats).Error; err != nil {
		return
	}

	queryCount, err := getDAppQuery(storage.DB, network, addresses, period)
	if err != nil {
		return
	}

	err = queryCount.Group("source").Count(&stats.Users).Error
	return
}

func getDAppQuery(db *gorm.DB, network string, addresses []string, period string) (*gorm.DB, error) {
	query := db.Table(models.DocOperations).
		Where("network = ?", network).
		Where("status = ?", constants.Applied)

	if len(addresses) > 0 {
		subQuery := db.Where("destination = ?", addresses[0])
		for i := 1; i < len(addresses); i++ {
			subQuery.Or("destination = ?", addresses[i])
		}
		query.Where(subQuery)
	}

	err := periodToRange(query, period)
	return query, err
}

func periodToRange(query *gorm.DB, period string) error {
	now := time.Now().UTC()
	switch period {
	case "year":
		now = now.AddDate(-1, 0, 0)
	case "month":
		now = now.AddDate(0, -1, 0)
	case "week":
		now = now.AddDate(0, 0, -7)
	case "day":
		now = now.AddDate(0, 0, -1)
	case "all":
		return nil
	default:
		return errors.Errorf("Unknown period value: %s", period)
	}
	query.Where("timestamp > ?", now)
	return nil
}

func addOperationSorting(query *gorm.DB) {
	query.Order("level desc, counter desc, nonce desc")
}
