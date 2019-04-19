package obj

import (
	"database/sql"

	// power
	_ "github.com/go-sql-driver/mysql"
)

// Datastore ...
type Datastore interface {
	ListMonitor(int, int) []Monitor
	ListMonitorByPages(int) []Monitor
	DestroyMonitor(int)
	CheckOwner(int, int) bool
	ListMonitorData(int) []MonitorData
	StoreHTTPMonitor(*NewMonitor)
	StorePingMonitor(*NewMonitor)
	StoreKeywordMonitor(*NewMonitor)
	MonitorHTTP(int)
	MonitorPing(int)
	MonitorKeyword(int)
}

// DB ...
type DB struct {
	*sql.DB
}

// OpenDB ...
func OpenDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
