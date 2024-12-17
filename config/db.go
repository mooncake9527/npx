package config

import (
	"strings"

	"github.com/mooncake9527/npx/common/consts"
	"gorm.io/gorm/logger"
)

type DB struct {
	DSN            string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`                                        //连接参数
	Disable        bool   `mapstructure:"disable" json:"disable" yaml:"disable"`                            //是否启用 默认true
	Driver         string `mapstructure:"driver" json:"driver" yaml:"driver"`                               //数据库类型
	Prefix         string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                               //全局表前缀，单独定义TableName则不生效
	MaxIdleConn    int    `mapstructure:"max-idle-conn" json:"max-idle-conn" yaml:"max-idle-conn"`          // 空闲中的最大连接数
	MaxOpenConn    int    `mapstructure:"max-open-conn" json:"max-open-conn" yaml:"max-open-conn"`          // 打开到数据库的最大连接数
	MaxLifetime    int    `mapstructure:"max-lifetime" json:"max-lifetime" yaml:"max-lifetime"`             // 链接重置时间（分）
	LogMode        string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                         // Gorm日志级别： silent、error、warn、info
	IgnoreNotFound bool   `mapstructure:"ignore-not-found" json:"ignore-not-found" yaml:"ignore-not-found"` //忽略无记录错误
	SlowThreshold  int    `mapstructure:"slow-threshold" json:"slow-threshold" yaml:"slow-threshold"`       // 慢查询 毫秒 大于0有效
	DryRun         bool   `mapstructure:"dry-run" json:"dry-run" yaml:"dry-run"`                            //
}

type DBCfg struct {
	DSN            string        `mapstructure:"dsn" json:"dsn" yaml:"dsn"`                                        //连接参数
	Driver         string        `mapstructure:"driver" json:"driver" yaml:"driver"`                               //数据库类型
	Prefix         string        `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                               //全局表前缀，单独定义TableName则不生效
	Singular       bool          `mapstructure:"singular" json:"singular" yaml:"singular"`                         //是否开启全局禁用复数，true表示开启
	MaxIdleConns   int           `mapstructure:"max-idle-conn" json:"max-idle-conn" yaml:"max-idle-conn"`          // 空闲中的最大连接数
	MaxOpenConns   int           `mapstructure:"max-open-conn" json:"max-open-conn" yaml:"max-open-conn"`          // 打开到数据库的最大连接数
	MaxLifetime    int           `mapstructure:"max-lifetime" json:"max-lifetime" yaml:"max-lifetime"`             // 链接重置时间（分）
	LogMode        string        `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                         // Gorm日志级别： silent、error、warn、info
	SlowThreshold  int           `mapstructure:"slow-threshold" json:"slow-threshold" yaml:"slow-threshold"`       // 慢查询 毫秒 大于0有效
	IgnoreNotFound bool          `mapstructure:"ignore-not-found" json:"ignore-not-found" yaml:"ignore-not-found"` //忽略无记录错误
	DryRun         bool          `mapstructure:"dry-run" json:"dry-run" yaml:"dry-run"`                            //
	DBS            map[string]DB `mapstructure:"dbs" json:"dbs" yaml:"dbs"`                                        //配置多db
}

func (c *DBCfg) GetDriver(dbname string) string {
	if dbname == consts.DbDefault {
		return c.Driver
	}
	if db, ok := c.DBS[dbname]; ok {
		return db.Driver
	}
	return ""
}

func (c *DBCfg) GetDSN(dbname string) string {
	if dbname == consts.DbDefault {
		return c.DSN
	}
	if db, ok := c.DBS[dbname]; ok {
		return db.DSN
	}
	return ""
}

func GetLogMode(logMode string) logger.LogLevel {
	switch strings.ToLower(logMode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

func (c *DBCfg) GetMaxIdleConn() int {
	if c.MaxIdleConns < 1 {
		c.MaxIdleConns = 10
	}
	return c.MaxIdleConns
}

func (c *DBCfg) GetMaxOpenConn() int {
	if c.MaxOpenConns < 1 {
		c.MaxOpenConns = 30
	}
	return c.MaxOpenConns
}

func (c *DBCfg) GetMaxLifetime() int {
	if c.MaxLifetime < 1 {
		c.MaxLifetime = 120
	}
	return c.MaxLifetime
}
