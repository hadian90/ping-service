package obj

import (
	"fmt"
	"time"

	"github.com/hadian90/ping-service/helper"
	"github.com/sparrc/go-ping"
)

// StorePingMonitor add new obj to database
func (db *DB) StorePingMonitor(d *NewMonitor) {
	d.Type = "ping"
	db.storeMonitor(d)
}

// -------------------------- CODE FOR MONITOR START -------------------

// MonitorPing ...
func (db *DB) MonitorPing(interval int) {

	var pm Monitor
	// get targer request
	stmtOut, err := db.Prepare("SELECT * FROM " + table + " WHERE type = 'ping' AND time_interval = ? AND last_request < ?")
	helper.ErrorHandler(err)
	defer stmtOut.Close()

	subTime := time.Now().UTC().Add(-time.Minute * time.Duration(interval)).Format("2006-01-02 15:04:05")

	rows, err := stmtOut.Query(interval, subTime)
	helper.ErrorHandler(err)
	defer rows.Close()

	for rows.Next() {
		// collect data from database
		err := rows.Scan(
			&pm.ID, &pm.UserID, &pm.Name, &pm.Destination, &pm.Type, &pm.PagesID, &pm.TimeInterval, &pm.CreatedAt, &pm.LastRequest)
		helper.ErrorHandler(err)

		// update target request
		stmtUpd, err := db.Prepare("UPDATE " + table + " SET last_request = ? WHERE id = ? ")
		helper.ErrorHandler(err)

		defer stmtUpd.Close()

		_, err = stmtUpd.Exec(time.Now(), pm.ID)
		helper.ErrorHandler(err)

		// make ping request to destination
		go db.pingRequestHandler(pm)
	}
}

func (db *DB) pingRequestHandler(pm Monitor) {
	fmt.Println("Pinging ", pm.Destination)
	// make ping request to destination
	pinger, err := ping.NewPinger(pm.Destination)
	if err != nil {
		db.updateMonitor(pm.ID, err.Error())
		return
	}
	pinger.Count = 1

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		// store condition on monitor data
		db.updateMonitor(pm.ID, fmt.Sprintf("OK with %v%% packet loss", stats.PacketLoss))
	}

	pinger.Run()
}

// -------------------------- CODE FOR MONITOR END -------------------
