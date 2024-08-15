package graph

import (
	"database/sql"
	"sync"

	"github.com/codecrafter404/bubble/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db              *sql.DB
	EventChannel    []chan *model.UpdateEvent
	EventChannelMux sync.RWMutex
}
