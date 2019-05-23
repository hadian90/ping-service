package obj

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

	var km KeywordMonitor
	// get targer request
	stmtOut, err := db.Prepare(
		"SELECT r.*, k.text, k.condition_type FROM " + table +
			" r, keyword k  WHERE r.type = 'keyword' AND k.request_id = r.id AND r.time_interval = ? AND r.last_request < ?")
	helper.ErrorHandler(err)
	defer stmtOut.Close()

	subTime := time.Now().UTC().Add(-time.Minute * time.Duration(interval)).Format("2006-01-02 15:04:05")

	rows, err := stmtOut.Query(interval, subTime)
	helper.ErrorHandler(err)
	defer rows.Close()

	for rows.Next() {
		// collect data from database
		err := rows.Scan(
			&km.ID, &km.UserID, &km.Name, &km.Destination, &km.Type,
			&km.PagesID, &km.TimeInterval, &km.CreatedAt, &km.LastRequest,
			&km.Text, &km.ConditionType)
		helper.ErrorHandler(err)

		// update target request
		stmtUpd, err := db.Prepare("UPDATE " + table + " SET last_request = ? WHERE id = ? ")
		helper.ErrorHandler(err)

		defer stmtUpd.Close()

		_, err = stmtUpd.Exec(time.Now(), km.ID)
		helper.ErrorHandler(err)

		go db.keywordRequestHandler(km)

	}
}

func (db *DB) keywordRequestHandler(km KeywordMonitor) {
	fmt.Println("Request ", km.Destination)
	// make http request to destination
	resp, err := http.Get(km.Destination)
	if err != nil {
		db.updateMonitor(km.ID, "Http", resp.StatusCode, err.Error())
		return
	}

	defer resp.Body.Close()

	fmt.Println("\n", km.Destination, " - HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	// scrap the website for text
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.ErrorHandler(err)
	}

	re := regexp.MustCompile("^(.*?(\b" + km.Text + "\b)[^$]*)$")
	comments := re.FindAllString(string(body), -1)
	if comments == nil {
		fmt.Println("No matches.")
		if km.ConditionType == "Found" {
			db.updateMonitor(km.ID, "Keyword", 404, "Keyword Not Found")
		} else {
			db.updateMonitor(km.ID, "Keyword", 200, "Keyword Found")
		}
	} else {
		for _, comment := range comments {
			fmt.Println(comment)
		}
		if km.ConditionType == "Found" {
			db.updateMonitor(km.ID, "Keyword", 200, "Keyword Found")
		} else {
			db.updateMonitor(km.ID, "Keyword", 404, "Keyword Not Found")
		}
	}
}

// -------------------------- CODE FOR MONITOR END -------------------
