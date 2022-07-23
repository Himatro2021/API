package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Himatro2021/API/internal/model"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_Create(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	repo := sessionRepo{
		db: kit.db,
	}

	session := &model.Session{
		ID:                    utils.GenerateID(),
		UserID:                utils.GenerateID(),
		AccessToken:           "at",
		RefreshToken:          "rt",
		AccessTokenExpiredAt:  time.Now(),
		RefreshTokenExpiredAt: time.Now().Add(1 * time.Hour),
	}

	t.Run("ok - created", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "sessions"`).
			WithArgs(
				session.UserID,
				session.AccessToken,
				session.RefreshToken,
				session.AccessTokenExpiredAt,
				session.RefreshTokenExpiredAt,
				session.ID,
			).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(session.ID))
		mock.ExpectCommit()

		err := repo.Create(ctx, session)

		assert.NoError(t, err)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "sessions"`).
			WithArgs(
				session.UserID,
				session.AccessToken,
				session.RefreshToken,
				session.AccessTokenExpiredAt,
				session.RefreshTokenExpiredAt,
				session.ID,
			).WillReturnError(errors.New("err db"))

		err := repo.Create(ctx, session)

		assert.Error(t, err)
	})
}
