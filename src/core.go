package src

import (
	"time"

	"github.com/hadian90/ping-service/obj"
)

// Core ...
type Core struct {
	obj.Datastore
}

// HTTPBackgroundService ...
func (c *Core) HTTPBackgroundService() {
	for range time.Tick(50 * time.Millisecond) {
		c.MonitorHTTP(1)
	}
}

// PingBackgroundService ...
func (c *Core) PingBackgroundService() {
	for range time.Tick(50 * time.Millisecond) {
		c.MonitorPing(1)
	}
}
