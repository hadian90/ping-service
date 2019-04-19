package obj

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hadian90/ping-service/helper"
)

// StoreHTTPMonitor add new obj to database
func (db *DB) StoreHTTPMonitor(d *NewMonitor) {
	d.Type = "http"
	db.storeMonitor(d)
}

// -------------------------- CODE FOR MONITOR START -------------------

// MonitorHTTP ...
func (db *DB) MonitorHTTP(interval int) {

	var md Monitor
	// get targer request
	stmtOut, err := db.Prepare("SELECT * FROM " + table + " WHERE type = 'http' AND time_interval = ? AND last_request < ?")
	helper.ErrorHandler(err)
	defer stmtOut.Close()

	subTime := time.Now().UTC().Add(-time.Minute * time.Duration(interval)).Format("2006-01-02 15:04:05")

	rows, err := stmtOut.Query(interval, subTime)
	helper.ErrorHandler(err)
	defer rows.Close()

	for rows.Next() {
		// collect data from database
		err := rows.Scan(
			&md.ID, &md.UserID, &md.Name, &md.Destination, &md.Type, &md.PagesID, &md.TimeInterval, &md.CreatedAt, &md.LastRequest)
		helper.ErrorHandler(err)

		// update target request
		stmtUpd, err := db.Prepare("UPDATE " + table + " SET last_request = ? WHERE id = ? ")
		helper.ErrorHandler(err)

		defer stmtUpd.Close()

		_, err = stmtUpd.Exec(time.Now(), md.ID)
		helper.ErrorHandler(err)

		// make http request to destination
		go db.httpRequestHandler(md)

	}
}

func (db *DB) httpRequestHandler(hm Monitor) {
	fmt.Println("Request ", hm.Destination)
	// make http request to destination
	resp, err := http.Get(hm.Destination)
	if err != nil {
		db.updateMonitor(hm.ID, err.Error())
		return
	}

	defer resp.Body.Close()

	fmt.Println("\n", hm.Destination, " - HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	// store response status
	db.updateMonitor(hm.ID, strconv.Itoa(resp.StatusCode))
}

// -------------------------- CODE FOR MONITOR END -------------------
