package rest

import (
	"github.com/Himatro2021/API/internal/model"
	"github.com/labstack/echo/v4"
)

// Service :nodoc:
type Service struct {
	group       *echo.Group
	userUsecase model.UserUsecase
}

// InitService self explained
func InitService(group *echo.Group, userUsecase model.UserUsecase) {
	service := Service{
		group:       group,
		userUsecase: userUsecase,
	}

	service.initRoutes()
}

func (s *Service) initRoutes() {
	s.group.POST("/members/invitations", s.handleCreateMemberInvitation())
}
