package monitor

import (
	"strconv"

	"github.com/hadian90/ping-service/helper"
	"github.com/hadian90/ping-service/obj"
	"github.com/labstack/echo"
)

// AddNewKeywordMonitor ...
func (m *Monitor) AddNewKeywordMonitor(ctx echo.Context) error {

	// bind data
	d := new(obj.NewMonitor)
	err := ctx.Bind(d)
	helper.ErrorHandler(err)

	// capture user id
	userid, err := strconv.Atoi(ctx.Request().Header.Get("user_id"))
	helper.ErrorHandler(err)

	d.UserID = userid

	m.StoreKeywordMonitor(d)

	return ctx.JSON(200, obj.Response{
		Success: true,
	})
}

// DeleteKeywordMonitor ...
func (m *Monitor) DeleteKeywordMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	m.DestroyMonitor(id)

	return ctx.JSON(200, obj.Response{
		Success: true,
		Message: "Ping request deleted",
	})
}

// DataKeywordMonitor ...
func (m *Monitor) DataKeywordMonitor(ctx echo.Context) error {

	id, err := strconv.Atoi(ctx.Param("id"))
	helper.ErrorHandler(err)

	return ctx.JSON(200, m.ListMonitorData(id))
}
