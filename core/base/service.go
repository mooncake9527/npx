package base

import (
	"github.com/mooncake9527/npx/core"
	"github.com/mooncake9527/npx/core/cache"
	"github.com/mooncake9527/x/xerrors/xerror"
	"gorm.io/gorm"
)

func NewDao(dbname string) *BaseDao {
	return &BaseDao{
		DbName: dbname,
	}
}

type BaseDao struct {
	DbName string
}

/*
* 获取数据库
 */
func (s *BaseDao) DB() *gorm.DB {
	return core.Db(s.DbName)
}

/*
* 获取缓存
 */
func (s *BaseDao) Cache() cache.ICache {
	return core.Cache
}

/*
* 创建 结构体model
 */
func (s *BaseDao) Create(model any) error {
	if err := s.DB().Create(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) CreateTx(tx *gorm.DB, model any) error {
	if err := tx.Create(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 更新整个模型 结构体model 注意空值
 */
func (s *BaseDao) Save(model any) error {
	if err := s.DB().Save(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) SaveTx(tx *gorm.DB, model any) error {
	if err := tx.Save(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 条件跟新
 */
func (s *BaseDao) UpdateWhere(model any, where any, updates map[string]any) error {
	if err := s.DB().Model(model).Where(where).Updates(updates).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) UpdateWhereTx(tx *gorm.DB, model any, where any, updates map[string]any) error {
	if err := tx.Model(model).Where(where).Updates(updates).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 模型更新
 */
func (s *BaseDao) UpdateWhereModel(where any, updates any) error {
	if err := s.DB().Where(where).Updates(updates).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) UpdateWhereModelTx(tx *gorm.DB, where any, updates any) error {
	if err := tx.Where(where).Updates(updates).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 根据模型id更新
 */
func (s *BaseDao) UpdateById(model any) error {
	if err := s.DB().Updates(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) UpdateByIdTx(tx *gorm.DB, model any) error {
	if err := tx.Updates(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

// DelModel model id 不能为空
func (s *BaseDao) DelModel(model any) error {
	if err := s.DB().Delete(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) DelWhereTx(tx *gorm.DB, model any) error {
	if err := tx.Delete(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 条件删除，模型 where 为map
 */
func (s *BaseDao) DelWhereMap(model any, where map[string]any) error {
	if err := s.DB().Where(where).Delete(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) DelWhereMapTx(tx *gorm.DB, model any, where map[string]any) error {
	if err := tx.Where(where).Delete(model).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* id 可以是uint64  也可以是[]uint64
 */

func (s *BaseDao) DelIds(model any, ids any) error {
	if err := s.DB().Delete(model, ids).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

func (s *BaseDao) DelIdsTx(tx *gorm.DB, model any, ids any) error {
	if err := tx.Delete(model, ids).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 根据id获取模型
 */
func (s *BaseDao) Get(id any, model any) error {
	if err := s.DB().First(model, id).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/**
* 条件查询
* where: where 查询条件model
* models: 代表查询返回的model数组
 */
func (s *BaseDao) GetByWhere(where any, models any) error {
	if err := s.DB().Where(where).Find(models).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/**
* 列表条件查询
* where: 条件查询
* models: 代表查询返回的model数组
 */
func (s *BaseDao) GetByMap(where map[string]any, models any) error {
	if err := s.DB().Where(where).Find(models).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/**
* 条数查询
* model: 查询条件
* count: 查询条数
 */
func (s *BaseDao) Count(model any, count *int64) error {
	if err := s.DB().Model(model).Where(model).Count(count).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/**
* 条数查询
* model: 查询条件
* count: 查询条数
 */
func (s *BaseDao) CountByMap(where map[string]any, model any, count *int64) error {
	if err := s.DB().Model(model).Where(where).Count(count).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/**
*	查询
* where 实现Query接口
 */
func (s *BaseDao) Query(where Query, models any) error {
	if err := s.DB().Scopes(s.MakeCondition(where)).Find(models).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 分页获取
 */
func (s *BaseDao) Page(where any, data any, total *int64, limit, offset int) error {
	if err := s.DB().Where(where).Limit(limit).Offset(offset).
		Find(data).Limit(-1).Offset(-1).Count(total).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 分页获取
 */
func (s *BaseDao) QueryPage(where Query, models any, total *int64, limit, offset int) error {
	if err := s.DB().Scopes(s.MakeCondition(where)).Limit(limit).Offset(offset).
		Find(models).Limit(-1).Offset(-1).Count(total).Error; err != nil {
		return xerror.New(err.Error())
	}
	return nil
}

/*
* 分页组装
 */
func (s *BaseDao) Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

/**
* chunk 查询
 */
func (s *BaseDao) Chunk(db *gorm.DB, size int, callback func(records []map[string]interface{}) error) error {
	var offset int
	for {
		var records []map[string]interface{}
		// 检索 size 条记录
		if err := db.Limit(size).Offset(offset).Find(&records).Error; err != nil {
			return xerror.New(err.Error())
		}
		// 如果没有更多记录，则退出循环
		if len(records) == 0 {
			break
		}
		// 调用回调函数处理记录
		if err := callback(records); err != nil {
			return err
		}
		// 更新偏移量
		offset += size
	}
	return nil
}

/**
* 查询条件组装
 */
func (s *BaseDao) MakeCondition(q Query) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &GormCondition{
			GormPublic: GormPublic{},
			Join:       make([]*GormJoin, 0),
		}
		ResolveSearchQuery(core.Cfg.DBCfg.GetDriver(s.DbName), q, condition, q.TableName())
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Joins(join.JoinOn)
			for k, v := range join.Where {
				db = db.Where(k, v...)
			}
			for k, v := range join.Or {
				db = db.Or(k, v...)
			}
			for _, o := range join.Order {
				db = db.Order(o)
			}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		return db
	}
}
