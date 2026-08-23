package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type Constant struct {
	Database DBConst    `json:"database"`
	Redis    RedisConst `json:"redis"`
}

type DBConst struct {
	MaxOpenDbConn int32         `json:"max_open_db_conn"`
	MaxIdleDbConn time.Duration `json:"max_idle_db_conn"`
	MaxDbLifeTime time.Duration `json:"max_db_life_time"`
}

type RedisConst struct {
}



func NewConstant() *Constant {
	return &Constant{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		Redis: RedisConst{},
	}
}

func (c *Constant) String() string {
	if c == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Constant{Database:%+v, Redis:%+v}",
			c.Database, c.Redis)
	}

	return string(data)
}
