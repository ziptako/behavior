package svc

import (
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/ziptako/behavior/db/model"
	"github.com/ziptako/behavior/internal/config"
)

type ServiceContext struct {
	Config         config.Config
	BehaviorsModel model.BehaviorsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:         c,
		BehaviorsModel: model.NewBehaviorsModel(conn, c.Cache),
	}
}
