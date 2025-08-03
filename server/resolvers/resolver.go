//go:generate go run generate.go
package resolvers

import (
	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/graph"
	"gorm.io/gorm"
)

type Resolver struct {
	EventChannel []chan *graph.UpdateEvent
	Db           *gorm.DB
	OrderChannel chan graph.Order
	Config       config.Config
}
