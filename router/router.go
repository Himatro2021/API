package router

import (
	"himatro-api/handler"
	"himatro-api/middleware"

	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func Router() *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.Logger())

	e.GET("/", handler.HomeGet)
	e.POST("/login", handler.Login)

	e.GET("/absensi/:absentID", handler.CheckAbsentForm)
	e.POST("/absensi/:absentID", handler.FillAbsentForm)
	e.PATCH("/absensi/:absentID", handler.UpdateAbsentListByAttendant)

	e.GET("/absensi/:absentID/result", handler.GetAbsentResult)

	e.GET("/admin", handler.Admin)
	e.POST("/admin/absensi", handler.InitAbsent, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/title", handler.UpdateFormTitle, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/participant", handler.UpdateFormParticipant, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/startAt", handler.UpdateFormStartAt, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/finishAt", handler.UpdateAbsentFormFinishAt, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/attendanceImageProof", handler.UpdateAbsentFormAttendanceImageProof, middleware.RequireLogin)
	e.PATCH("/admin/absensi/:absentID/execuseImageProof", handler.UpdateAbsentFormExecuseImageProof, middleware.RequireLogin)

	return e
}
