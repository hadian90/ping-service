package monitor

import (
	"strconv"

	"github.com/hadian90/ping-service/helper"
	"github.com/hadian90/ping-service/obj"
	"github.com/labstack/echo"
)

// AddNewHTTPMonitor ...
func (m *Monitor) AddNewHTTPMonitor(ctx echo.Context) error {

	// bind data
	d := new(obj.NewMonitor)
	err := ctx.Bind(d)
	helper.ErrorHandler(err)

	// capture user id
	userid, err := strconv.Atoi(ctx.Request().Header.Get("user_id"))
	helper.ErrorHandler(err)

	d.UserID = userid

	m.StoreHTTPMonitor(d)

	return ctx.JSON(200, obj.Response{
		Success: true,
		Message: "New HTTP request added",
	})
}

// DeleteHTTPMonitor ...
func (m *Monitor) DeleteHTTPMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	// capture user id
	userid, err := strconv.Atoi(ctx.Request().Header.Get("user_id"))
	helper.ErrorHandler(err)

	if !m.CheckOwner(id, userid) {
		return ctx.JSON(403, obj.Response{
			Success: false,
			Message: "Wrong permission",
		})
	}

	m.DestroyMonitor(id)

	return ctx.JSON(200, obj.Response{
		Success: true,
		Message: "HTTP request deleted",
	})
}

// DataHTTPMonitor ...
func (m *Monitor) DataHTTPMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	// capture user id
	userid, err := strconv.Atoi(ctx.Request().Header.Get("user_id"))
	helper.ErrorHandler(err)

	if !m.CheckOwner(id, userid) {
		return ctx.JSON(403, obj.Response{
			Success: false,
			Message: "Wrong permission",
		})
	}

	return ctx.JSON(200, m.ListMonitorData(id))
}
