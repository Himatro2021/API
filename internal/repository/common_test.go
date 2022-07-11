package repository

import (
	"log"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Himatro2021/API/internal/model/mock"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadConf :nodoc:
func LoadConf() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}
}

type repoTestKit struct {
	dbmock       sqlmock.Sqlmock
	db           *gorm.DB
	mockUserRepo *mock.MockUserRepository
	ctrl         *gomock.Controller
}

func initializeRepoTestKit(t *testing.T) *repoTestKit {
	ctrl := gomock.NewController(t)
	userRepository := mock.NewMockUserRepository(ctrl)

	dbconn, dbmock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: dbconn}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return &repoTestKit{
		ctrl:         ctrl,
		mockUserRepo: userRepository,
		db:           gormDB,
		dbmock:       dbmock,
	}
}
