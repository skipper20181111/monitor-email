package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DB struct {
		DataSource string
	}
	ServerInfo struct {
		Url string
	}
	CDP struct {
		ServicesUrl string
	}
	Monitor struct {
		TimeGapMinute  int
		HistoryKeepDay int
	}
}
