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
		c.MonitorHTTP(5)
		c.MonitorHTTP(10)
		c.MonitorHTTP(15)
	}
}

// PingBackgroundService ...
func (c *Core) PingBackgroundService() {
	for range time.Tick(50 * time.Millisecond) {
		c.MonitorPing(5)
		c.MonitorPing(10)
		c.MonitorPing(15)
	}
}

// KeywordBackgroundService ...
func (c *Core) KeywordBackgroundService() {
	for range time.Tick(50 * time.Millisecond) {
		c.MonitorKeyword(5)
		c.MonitorKeyword(10)
		c.MonitorKeyword(15)
	}
}
