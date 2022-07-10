package rest

import (
	"github.com/Himatro2021/API/internal/model"
	"github.com/labstack/echo/v4"
)

type RESTService struct {
	group       *echo.Group
	userUsecase model.UserUsecase
}

func InitRESTService(group *echo.Group, userUsecase model.UserUsecase) {
	service := RESTService{
		group:       group,
		userUsecase: userUsecase,
	}

	service.initRoutes()
}

func (s *RESTService) initRoutes() {
	s.group.POST("/members/invitations", s.handleCreateMemberInvitation())
}
