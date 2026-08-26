package db

import (
	"GoSingleBoot/internal/config"
	"context"
	"database/sql"
	"math/rand/v2"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresClient struct {
	Master *bun.DB
	slaves []*bun.DB
}

// Slave 如果没有就默认使用Master避免报错，如果只有1个，直接使用，1个以上就随机
func (c *PostgresClient) Slave() *bun.DB {
	i := len(c.slaves)
	if i == 0 {
		return c.Master
	} else if i == 1 {
		return c.slaves[0]
	} else {
		r := rand.IntN(i - 1)
		return c.slaves[r]
	}
}

var Client *PostgresClient

func NewPostgresClient() {
	masterDsn := config.Config.DatabaseCfg.Master
	slavesDsn := config.Config.DatabaseCfg.Slaves
	var client PostgresClient

	masterClient, err := Connect(masterDsn)
	if err != nil {
		panic("Master 数据库 " + masterDsn + " 连接失败 :" + err.Error())
	}
	client.Master = masterClient

	if len(slavesDsn) > 0 {
		slavesList := make([]*bun.DB, 0)
		for _, slaveDsn := range slavesDsn {
			slaveClient, err := Connect(slaveDsn)
			if err != nil {
				panic("Slave 数据库 " + slaveDsn + " 连接失败 :" + err.Error())
			}
			slavesList = append(slavesList, slaveClient)
		}
		client.slaves = slavesList
	}

	Client = &client
}

func Connect(dsn string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
