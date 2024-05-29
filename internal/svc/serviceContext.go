package svc

import (
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"monitor/cachemodel"
	"monitor/internal/config"
	"time"
)

const (
	localCacheExpire = time.Duration(time.Minute * 20)
	SystemListKey    = "SystemListKey"
	Keystr           = "W3WxhhoA9E9VIteCYbnhUTxDbtk2nP1Z"
	EmailListKey     = "EmailListKey"
)

type ServiceContext struct {
	Config     config.Config
	LocalCache *collection.Cache
	Monitor    cachemodel.ShrcbMonitorModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		LocalCache: localCache,
		Monitor:    cachemodel.NewShrcbMonitorModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
