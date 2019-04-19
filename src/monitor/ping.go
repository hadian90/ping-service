package monitor

import (
	"strconv"

	"github.com/hadian90/ping-service/helper"
	"github.com/hadian90/ping-service/obj"
	"github.com/labstack/echo"
)

// AddNewPingMonitor ...
func (m *Monitor) AddNewPingMonitor(ctx echo.Context) error {

	// bind data
	d := new(obj.NewMonitor)
	err := ctx.Bind(d)
	helper.ErrorHandler(err)

	m.StorePingMonitor(d)

	return ctx.JSON(200, obj.Response{
		Success: true,
		Message: "Ping request added",
	})
}

// DeletePingMonitor ...
func (m *Monitor) DeletePingMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	m.DestroyMonitor(id)

	return ctx.JSON(200, obj.Response{
		Success: true,
		Message: "Ping request deleted",
	})
}

// DataPingMonitor ...
func (m *Monitor) DataPingMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	return ctx.JSON(200, m.ListMonitorData(id))
}
