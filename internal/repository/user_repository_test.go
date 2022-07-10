package repository

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateInvitation(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	name := "lucky"
	email := "lucky@test.ting"

	t.Run("ok - created", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "user_invitations"`).
			WithArgs(email, name, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(utils.GenerateID()))
		mock.ExpectCommit()

		_, err := repo.CreateInvitation(ctx, name, email)

		assert.NoError(t, err)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "user_invitations"`).
			WithArgs(email, name, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		_, err := repo.CreateInvitation(ctx, name, email)

		assert.Error(t, err)
	})
}
