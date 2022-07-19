package rest

import (
	"github.com/Himatro2021/API/internal/model"
	"github.com/labstack/echo/v4"
)

// Service :nodoc:
type Service struct {
	group         *echo.Group
	userUsecase   model.UserUsecase
	absentUsecase model.AbsentUsecase
}

// InitService self explained
func InitService(group *echo.Group, userUsecase model.UserUsecase, absentUsecase model.AbsentUsecase) {
	service := Service{
		group:         group,
		userUsecase:   userUsecase,
		absentUsecase: absentUsecase,
	}

	service.initRoutes()
}

func (s *Service) initRoutes() {
	s.initMemberHandlerRoutes()
	s.initAbsentHandlerRoutes()
}

func (s *Service) initMemberHandlerRoutes() {
	s.group.POST("/members/invitations", s.handleCreateMemberInvitation())
}

func (s *Service) initAbsentHandlerRoutes() {
	s.group.GET("/absent/form/", s.handleGetAllAbsentForms())
	s.group.GET("/absent/form/:id/", s.handleGetFormByID())
	s.group.GET("/absent/form/:id/result/", s.handleGetParticipantsByFormID())

	s.group.POST("/absent/form/", s.handleCreateAbsentForm())
	s.group.POST("/absent/form/:id/", s.handleFillAbsentFormByAttendee())

	s.group.PUT("/absent/form/:id/", s.handleUpdateAbsentForm())

	s.group.PATCH("/absent/form/attendance/:id/", s.handleUpdateAbsentListByAttendee())
}
