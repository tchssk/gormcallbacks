package v1

import (
	"sync"

	"github.com/jinzhu/gorm"
)

type Config struct {
	CreateCount map[string]int
	QueryCount  map[string]int
	DeleteCount map[string]int
	UpdateCount map[string]int
}

func Count(m map[string]int) func(*gorm.Scope) {
	var mu sync.Mutex
	return func(scope *gorm.Scope) {
		mu.Lock()
		m[scope.TableName()]++
		mu.Unlock()
	}
}

func RegisterCountCallbacks(db *gorm.DB, config *Config) {
	createCallback := db.Callback().Create()
	if config.CreateCount != nil {
		createCallback.After("gorm:create").Register("count", Count(config.CreateCount))
	}

	queryCallback := db.Callback().Query()
	if config.QueryCount != nil {
		queryCallback.After("query").Register("count", Count(config.QueryCount))
	}

	deleteCallback := db.Callback().Delete()
	if config.DeleteCount != nil {
		deleteCallback.After("gorm:delete").Register("count", Count(config.DeleteCount))
	}

	updateCallback := db.Callback().Update()
	if config.UpdateCount != nil {
		updateCallback.After("gorm:update").Register("count", Count(config.UpdateCount))
	}
}
