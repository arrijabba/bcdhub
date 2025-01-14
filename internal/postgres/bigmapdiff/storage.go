package bigmapdiff

import (
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/bigmapdiff"
	"github.com/baking-bad/bcdhub/internal/postgres/consts"
	"github.com/baking-bad/bcdhub/internal/postgres/core"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Storage -
type Storage struct {
	*core.Postgres
}

// NewStorage -
func NewStorage(pg *core.Postgres) *Storage {
	return &Storage{pg}
}

func bigMapKey(network, keyHash string, ptr int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("network = ?", network).
			Where("key_hash = ?", keyHash).
			Where("ptr = ?", ptr)
	}
}

// CurrentByKey -
func (storage *Storage) Current(network, keyHash string, ptr int64) (data bigmapdiff.BigMapState, err error) {
	if ptr < 0 {
		err = errors.Wrapf(consts.ErrInvalidPointer, "%d", ptr)
		return
	}

	err = storage.DB.Table(models.DocBigMapState).
		Scopes(bigMapKey(network, keyHash, ptr)).
		First(&data).
		Error

	return
}

// GetForAddress -
func (storage *Storage) GetForAddress(address string) (response []bigmapdiff.BigMapDiff, err error) {
	err = storage.DB.Table(models.DocBigMapDiff).
		Scopes(core.Contract(address)).
		Group("key_hash").
		Order("level desc").
		Find(&response).
		Error
	return
}

// GetByAddress -
func (storage *Storage) GetByAddress(network, address string) (response []bigmapdiff.BigMapDiff, err error) {
	err = storage.DB.Table(models.DocBigMapDiff).
		Scopes(core.NetworkAndContract(network, address)).
		Order("level desc").
		Find(&response).
		Error
	return
}

// GetValuesByKey -
func (storage *Storage) GetValuesByKey(keyHash string) (response []bigmapdiff.BigMapDiff, err error) {
	err = storage.DB.Table(models.DocBigMapDiff).
		Where("key_hash = ?", keyHash).
		Group("network, contract, ptr").
		Order("level desc").
		Find(&response).
		Error
	return
}

// Count -
func (storage *Storage) Count(network string, ptr int64) (count int64, err error) {
	err = storage.DB.Table(models.DocBigMapState).
		Where("network = ?", network).
		Where("ptr = ?", ptr).
		Count(&count).
		Error
	return
}

// Previous -
func (storage *Storage) Previous(filters []bigmapdiff.BigMapDiff, lastID int64, address string) (response []bigmapdiff.BigMapDiff, err error) {
	query := storage.DB.Table(models.DocBigMapDiff).
		Scopes(core.Contract(address)).
		Select("MAX(id) as id")

	if lastID > 0 {
		query.Where("id < ?", lastID)
	}
	if len(filters) > 0 {
		tx := storage.DB.Where(
			storage.DB.Where("key_hash = ?", filters[0].KeyHash).Where("ptr = ? ", filters[0].Ptr),
		)
		for i := 1; i < len(filters); i++ {
			tx.Or(
				storage.DB.Where("key_hash = ?", filters[i].KeyHash).Where("ptr = ? ", filters[i].Ptr),
			)
		}
		query.Where(tx)
	}

	query.Group("key_hash")

	err = storage.DB.Table(models.DocBigMapDiff).
		Where("id IN (?)", query).
		Find(&response).Error

	return
}

// GetForOperation -
func (storage *Storage) GetForOperation(hash string, counter int64, nonce *int64) (response []*bigmapdiff.BigMapDiff, err error) {
	query := storage.DB.Table(models.DocBigMapDiff).
		Where("operation_hash = ?", hash).
		Where("operation_counter = ?", counter)

	if nonce == nil {
		query.Where("operation_nonce IS NULL")
	} else {
		query.Where("operation_nonce = ?", *nonce)
	}

	return response, query.Find(&response).Error
}

// GetUniqueForOperation -
func (storage *Storage) GetUniqueForOperation(hash string, counter int64, nonce *int64) (response []bigmapdiff.BigMapDiff, err error) {
	subQuery := storage.DB.Table(models.DocBigMapDiff).
		Select("MAX(id) as id").
		Where("operation_hash = ?", hash).
		Where("operation_counter = ?", counter)

	if nonce == nil {
		subQuery.Where("operation_nonce IS NULL")
	} else {
		subQuery.Where("operation_nonce = ?", *nonce)
	}

	subQuery.Group("key_hash, ptr")

	err = storage.DB.Table(models.DocBigMapDiff).
		Where("id IN (?)", subQuery).
		Find(&response).Error

	return
}

// GetByPtrAndKeyHash -
func (storage *Storage) GetByPtrAndKeyHash(ptr int64, network, keyHash string, size, offset int64) ([]bigmapdiff.BigMapDiff, int64, error) {
	if ptr < 0 {
		return nil, 0, errors.Wrapf(consts.ErrInvalidPointer, "%d", ptr)
	}
	limit := storage.GetPageSize(size)

	var response []bigmapdiff.BigMapDiff
	if err := storage.DB.
		Scopes(core.Network(network), core.OrderByLevelDesc).
		Where("key_hash = ?", keyHash).
		Where("ptr = ?", ptr).
		Limit(limit).
		Offset(int(offset)).
		Find(&response).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	err := storage.DB.Table(models.DocBigMapDiff).
		Scopes(core.Network(network)).
		Where("key_hash = ?", keyHash).
		Where("ptr = ?", ptr).
		Count(&count).Error

	return response, count, err
}

// GetByPtr -
func (storage *Storage) GetByPtr(address, network string, ptr int64) (response []bigmapdiff.BigMapDiff, err error) {
	subQuery := storage.DB.Table(models.DocBigMapDiff).
		Select("max(id) as id").
		Scopes(core.NetworkAndContract(network, address)).
		Where("ptr = ?", ptr).
		Group("key_hash").
		Order("id desc")

	query := storage.DB.Table(models.DocBigMapDiff).Joins("inner join (?) as bmd on bmd.id = big_map_diffs.id", subQuery)
	return response, query.Find(&response).Error
}

// Get -
func (storage *Storage) Get(ctx bigmapdiff.GetContext) ([]bigmapdiff.Bucket, error) {
	if *ctx.Ptr < 0 {
		return nil, errors.Errorf("Invalid pointer value: %d", *ctx.Ptr)
	}

	var bmd []bigmapdiff.Bucket
	subQuery := storage.buildGetContext(ctx)

	query := storage.DB.Table(models.DocBigMapDiff).Select("*, bmd.keys_count").Joins("inner join (?) as bmd on bmd.id = big_map_diffs.id", subQuery)
	err := query.Find(&bmd).Error
	return bmd, err
}

// GetByIDs -
func (storage *Storage) GetByIDs(ids ...int64) (result []bigmapdiff.BigMapDiff, err error) {
	err = storage.DB.Table(models.DocBigMapDiff).Order("id asc").Find(&result, ids).Error
	return
}

// GetStats -
func (storage *Storage) GetStats(network string, ptr int64) (stats bigmapdiff.Stats, err error) {
	totalQuery := storage.DB.Table(models.DocBigMapState).
		Select("count(contract) as total, contract").
		Where("network = ?", network).
		Where("ptr = ?", ptr).
		Group("contract")

	activeQuery := storage.DB.Table(models.DocBigMapState).
		Select("count(contract) as active, contract").
		Where("network = ?", network).
		Where("ptr = ?", ptr).
		Where("removed = false").
		Group("contract")

	err = storage.DB.
		Raw("select q1.total as total, q1.contract as contract, q2.active as active from ((?) q1 join (?) q2 on q1.contract = q2.contract)", totalQuery, activeQuery).
		Scan(&stats).
		Error

	return
}

// CurrentByContract -
func (storage *Storage) CurrentByContract(network, contract string) (keys []bigmapdiff.BigMapState, err error) {
	err = storage.DB.Table(models.DocBigMapState).
		Where("network = ?", network).
		Where("contract = ?", contract).
		Find(&keys).
		Error

	return
}

// StatesChangedAfter -
func (storage *Storage) StatesChangedAfter(network string, level int64) (states []bigmapdiff.BigMapState, err error) {
	err = storage.DB.Table(models.DocBigMapState).
		Where("network = ?", network).
		Where("last_update_level > ?", level).
		Find(&states).
		Error
	return
}

// LastDiff -
func (storage *Storage) LastDiff(network string, ptr int64, keyHash string, skipRemoved bool) (diff bigmapdiff.BigMapDiff, err error) {
	query := storage.DB.Table(models.DocBigMapDiff).
		Scopes(bigMapKey(network, keyHash, ptr))

	if skipRemoved {
		query.Where("value is not null")
	}

	err = query.Order("id desc").Scan(&diff).Error
	return
}
