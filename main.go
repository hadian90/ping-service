package main

import (
	"database/sql"

	"github.com/hadian90/ping-service/config"
	"github.com/hadian90/ping-service/helper"
	"github.com/hadian90/ping-service/obj"
	"github.com/hadian90/ping-service/src"
	"github.com/hadian90/ping-service/src/monitor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	id             int
	destinationURL string
	webhookPath    string
	timeInterval   int
	lastRequest    sql.NullString
)

func main() {

	db, err := obj.OpenDB("hadian90:hadian90@/ping-service")
	helper.ErrorHandler(err)

	defer db.Close()

	// calling all collected class
	core := &src.Core{db}
	monitor := &monitor.Monitor{db}

	// run monitor service
	go core.HTTPBackgroundService()
	go core.PingBackgroundService()
	go core.KeywordBackgroundService()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/monitor")
	g.Use(config.UserMiddleware)
	// Routes
	g.GET("", monitor.GetMonitor)

	g.POST("/http/add", monitor.AddNewHTTPMonitor)
	g.GET("/http/:id/delete", monitor.DeleteHTTPMonitor)
	g.GET("/http/:id/data", monitor.DataHTTPMonitor)

	g.POST("/ping/add", monitor.AddNewPingMonitor)
	g.GET("/ping/:id/delete", monitor.DeletePingMonitor)
	g.GET("/ping/:id/data", monitor.DataPingMonitor)

	g.POST("/keyword/add", monitor.AddNewKeywordMonitor)
	g.GET("/keyword/:id/delete", monitor.DeleteKeywordMonitor)
	g.GET("/htkeywordtp/:id/data", monitor.DataKeywordMonitor)

	e.GET("/public/:id", monitor.GetMonitorByPages)

	// Start server
	e.Logger.Fatal(e.Start(":1337"))

}
