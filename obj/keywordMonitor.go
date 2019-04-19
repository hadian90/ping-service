package obj

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hadian90/ping-service/helper"
)

// StoreKeywordMonitor add new obj to database
func (db *DB) StoreKeywordMonitor(d *NewMonitor) {
	d.Type = "keyword"
	db.storeMonitor(d)
}

// -------------------------- CODE FOR MONITOR START -------------------

// MonitorKeyword ...
func (db *DB) MonitorKeyword(interval int) {

	var km Monitor
	// get targer request
	stmtOut, err := db.Prepare("SELECT * FROM " + table + " WHERE type = 'keyword' AND time_interval = ? AND last_request < ?")
	helper.ErrorHandler(err)
	defer stmtOut.Close()

	subTime := time.Now().UTC().Add(-time.Minute * time.Duration(interval)).Format("2006-01-02 15:04:05")

	rows, err := stmtOut.Query(interval, subTime)
	helper.ErrorHandler(err)
	defer rows.Close()

	for rows.Next() {
		// collect data from database
		err := rows.Scan(
			&km.ID, &km.UserID, &km.Name, &km.Destination, &km.Type, &km.PagesID, &km.TimeInterval, &km.CreatedAt, &km.LastRequest)
		helper.ErrorHandler(err)

		// update target request
		stmtUpd, err := db.Prepare("UPDATE " + table + " SET last_request = ? WHERE id = ? ")
		helper.ErrorHandler(err)

		defer stmtUpd.Close()

		_, err = stmtUpd.Exec(time.Now(), km.ID)
		helper.ErrorHandler(err)

		// make http request to destination
		resp, err := http.Get(km.Destination)
		if err != nil {
			// if error skip loop
			break
		}
		defer resp.Body.Close()

		fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

		// store response status
		dataStmtIns, err := db.Prepare("INSERT INTO monitor_data (monitor_type, request_id, status) VALUES ( ?, ?, ? )")
		helper.ErrorHandler(err)

		defer dataStmtIns.Close()

		_, err = dataStmtIns.Exec("Keyword", km.ID, resp.StatusCode)
		helper.ErrorHandler(err)

	}
}

// -------------------------- CODE FOR MONITOR END -------------------
