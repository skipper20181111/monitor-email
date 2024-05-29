package cachemodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShrcbMonitorModel = (*customShrcbMonitorModel)(nil)

type (
	// ShrcbMonitorModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShrcbMonitorModel.
	ShrcbMonitorModel interface {
		shrcbMonitorModel
	}

	customShrcbMonitorModel struct {
		*defaultShrcbMonitorModel
	}
)

// NewShrcbMonitorModel returns a model for the database table.
func NewShrcbMonitorModel(conn sqlx.SqlConn) ShrcbMonitorModel {
	return &customShrcbMonitorModel{
		defaultShrcbMonitorModel: newShrcbMonitorModel(conn),
	}
}
