package rbac

import (
	"log"
	"sync"

	"github.com/MatinHAB05/2pi/config"
	"github.com/MatinHAB05/2pi/internal/infrastructure/database"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type RBACEnforcer interface {
	GetEnforcer() *casbin.Enforcer
	Enforce(rvals ...any) (bool, error)
}

type Casbin struct {
	enforcer *casbin.Enforcer
}

var (
	casbinOnce     sync.Once
	casbinInstance *Casbin
)

func NewCasbin(casbinConfig *config.Casbin, gormDB database.Database) RBACEnforcer {
	casbinOnce.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(gormDB.GetDB())
		if err != nil {
			log.Fatalf("Failed to create Casbin adapter: %v", err)
		}

		m, err := model.NewModelFromFile(casbinConfig.ModelConfigFilePath)
		if err != nil {
			log.Fatalf("Failed to parse Casbin model configuration: %v", err)
		}

		enforcer, err := casbin.NewEnforcer(m, adapter)
		if err != nil {
			log.Fatalf("Failed to create Casbin enforcer: %v", err)
		}

		err = enforcer.LoadPolicy()
		if err != nil {
			log.Fatalf("Failed to load Casbin policies from DB: %v", err)
		}

		casbinInstance = &Casbin{enforcer: enforcer}
	})

	return casbinInstance
}

// GetEnforcer returns the underlying Casbin enforcer instance.
func (c *Casbin) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

// Enforce evaluates raw arguments against the Casbin model.
func (c *Casbin) Enforce(rvals ...any) (bool, error) {
	return c.enforcer.Enforce(rvals...)
}
