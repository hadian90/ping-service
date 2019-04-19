package obj

import (
	"time"

	"github.com/hadian90/ping-service/helper"
)

type (
	// NewMonitor binding object on store request
	NewMonitor struct {
		UserID       int    `json:"user_id"`
		Name         string `json:"name"`
		Destination  string `json:"destination"`
		Type         string `json:"type"`
		PagesID      int    `json:"pages_id"`
		TimeInterval int    `json:"time_interval"`
	}

	// Monitor ...
	Monitor struct {
		ID           int    `json:"id"`
		UserID       int    `json:"user_id"`
		Name         string `json:"name"`
		Destination  string `json:"destination"`
		Type         string `json:"type"`
		PagesID      int    `json:"pages_id"`
		TimeInterval int    `json:"time_interval"`
		CreatedAt    string `json:"created_at"`
		LastRequest  string `json:"last_request"`
		Status       string `json:"status"`
	}

	// MonitorData ...
	MonitorData struct {
		ID               int    `json:"id"`
		MonitorType      string `json:"monitor_type"`
		RequestTimestamp string `json:"request_timestamp"`
		RequestID        int    `json:"request_id"`
		Status           string `json:"status"`
	}
)

var table = "monitor_request"

var dataTable = "monitor_data"

// ListMonitor ...
func (db *DB) ListMonitor(userid int, pagesid int) []Monitor {
	// prepare statement
	httpStmt, err := db.Prepare(
		"SELECT * FROM " + table + " WHERE user_id = ? AND pages_id = ?")
	helper.ErrorHandler(err)

	defer httpStmt.Close()

	rows, err := httpStmt.Query(userid, pagesid)
	helper.ErrorHandler(err)

	resultArray := make([]Monitor, 0)

	for rows.Next() {
		var md Monitor

		err := rows.Scan(
			&md.ID, &md.UserID, &md.Name, &md.Destination, &md.Type, &md.PagesID, &md.TimeInterval, &md.CreatedAt, &md.LastRequest)
		helper.ErrorHandler(err)

		// get monitor data
		md.Status = db.latestMonitorStatus(md.ID)

		resultArray = append(resultArray, md)
	}

	return resultArray
}

// ListMonitorByPages ...
func (db *DB) ListMonitorByPages(pagesid int) []Monitor {
	// prepare statement
	httpStmt, err := db.Prepare(
		"SELECT * FROM " + table + " WHERE pages_id = ?")
	helper.ErrorHandler(err)

	defer httpStmt.Close()

	rows, err := httpStmt.Query(pagesid)
	helper.ErrorHandler(err)

	resultArray := make([]Monitor, 0)

	for rows.Next() {
		var md Monitor

		err := rows.Scan(
			&md.ID, &md.UserID, &md.Name, &md.Destination, &md.Type, &md.PagesID, &md.TimeInterval, &md.CreatedAt, &md.LastRequest)
		helper.ErrorHandler(err)

		// get monitor data
		md.Status = db.latestMonitorStatus(md.ID)

		resultArray = append(resultArray, md)
	}

	return resultArray
}

func (db *DB) storeMonitor(d *NewMonitor) {

	// prepare statement
	httpStmtIns, err := db.Prepare(
		"INSERT INTO " + table + " (user_id, name, destination, type, pages_id, time_interval, last_request) VALUES ( ?, ?, ?, ?, ?, ?, ?)")
	helper.ErrorHandler(err)

	defer httpStmtIns.Close()

	_, err = httpStmtIns.Exec(d.UserID, d.Name, d.Destination, d.Type, d.PagesID, d.TimeInterval, time.Now().UTC())
	helper.ErrorHandler(err)

}

// DestroyMonitor ...
func (db *DB) DestroyMonitor(id int) {
	// prepare statement
	httpStmt, err := db.Prepare("DELETE FROM " + table + " WHERE id = ?")
	helper.ErrorHandler(err)

	defer httpStmt.Close()

	_, err = httpStmt.Exec(id)
	helper.ErrorHandler(err)
}

// CheckOwner ..
func (db *DB) CheckOwner(id int, userID int) bool {
	// prepare statement
	httpStmt, err := db.Prepare(
		"SELECT COUNT(*) FROM " + table + " WHERE id = ? AND user_id = ?")
	helper.ErrorHandler(err)

	defer httpStmt.Close()

	rows, err := httpStmt.Query(id, userID)
	helper.ErrorHandler(err)

	var result int

	for rows.Next() {
		err := rows.Scan(&result)
		helper.ErrorHandler(err)
	}

	if result == 1 {
		return true
	}
	return false
}

// ListMonitorData ...
func (db *DB) ListMonitorData(id int) []MonitorData {
	// prepare statement
	httpStmt, err := db.Prepare(
		"SELECT * FROM monitor_data WHERE request_id = ? ORDER BY request_timestamp")
	helper.ErrorHandler(err)

	defer httpStmt.Close()

	rows, err := httpStmt.Query(id)
	helper.ErrorHandler(err)

	resultArray := make([]MonitorData, 0)

	for rows.Next() {
		var md MonitorData

		err := rows.Scan(&md.ID, &md.MonitorType, &md.RequestTimestamp, &md.RequestID, &md.Status)
		helper.ErrorHandler(err)

		resultArray = append(resultArray, md)
	}

	return resultArray
}

// -------------------------- CODE FOR MONITOR -------------------

func (db *DB) latestMonitorStatus(id int) string {
	dataStmt, err := db.Prepare("SELECT status FROM " + dataTable + " WHERE request_id = ? ORDER BY request_timestamp LIMIT 1")
	helper.ErrorHandler(err)

	defer dataStmt.Close()

	rows, err := dataStmt.Query(id)
	helper.ErrorHandler(err)

	var status = "Running"

	for rows.Next() {

		err := rows.Scan(&status)
		helper.ErrorHandler(err)

	}

	return status
}

func (db *DB) updateMonitor(id int, status string) {
	// store response status
	dataStmtIns, err := db.Prepare("INSERT INTO monitor_data (monitor_type, request_id, status) VALUES ( ?, ?, ? )")
	helper.ErrorHandler(err)

	defer dataStmtIns.Close()

	_, err = dataStmtIns.Exec("Http", id, status)
	helper.ErrorHandler(err)
}

// -------------------------- CODE FOR MONITOR -------------------
