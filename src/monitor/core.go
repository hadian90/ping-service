package monitor

import (
	"strconv"

	"github.com/hadian90/ping-service/helper"
	"github.com/hadian90/ping-service/obj"
	"github.com/labstack/echo"
)

// Monitor ...
type Monitor struct {
	obj.Datastore
}

// GetMonitor ...
func (m *Monitor) GetMonitor(ctx echo.Context) error {

	// capture user id
	userid, err := strconv.Atoi(ctx.Request().Header.Get("user_id"))
	helper.ErrorHandler(err)

	pagesid, err := strconv.Atoi(ctx.Request().Header.Get("pages_id"))
	helper.ErrorHandler(err)

	return ctx.JSON(200, m.ListMonitor(userid, pagesid))
}

// GetMonitorByPages ...
func (m *Monitor) GetMonitorByPages(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	return ctx.JSON(200, m.ListMonitorByPages(id))
}
