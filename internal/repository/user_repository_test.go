package repository

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_CreateInvitation(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	invitation := &model.UserInvitation{
		ID:     utils.GenerateID(),
		Email:  "test@ting.com",
		Name:   "ini nama",
		Role:   rbac.RoleAdmin,
		Status: model.InvitationStatusPending,
	}

	t.Run("ok - created", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "user_invitations"`).
			WithArgs(invitation.MailServiceID, invitation.Email, invitation.Name, invitation.InvitationCode, invitation.Role, invitation.Status, invitation.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(utils.GenerateID()))
		mock.ExpectCommit()

		err := repo.CreateInvitation(ctx, invitation)

		assert.NoError(t, err)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "user_invitations"`).
			WithArgs(invitation.MailServiceID, invitation.Email, invitation.Name, invitation.InvitationCode, invitation.Role, invitation.Status, invitation.ID).
			WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		err := repo.CreateInvitation(ctx, invitation)

		assert.Error(t, err)
	})
}

func TestUserRepository_IsEmailAlreadyInvited(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	email := "lucky@test.ting"

	t.Run("ok - email already exists", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "user_invitations"`).
			WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(utils.GenerateID()))

		isExists, err := repo.IsEmailAlreadyInvited(ctx, email)

		assert.NoError(t, err)
		assert.Equal(t, true, isExists)
	})

	t.Run("ok - email not exists", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "user_invitations"`).
			WithArgs(email).WillReturnError(gorm.ErrRecordNotFound)

		isExists, err := repo.IsEmailAlreadyInvited(ctx, email)

		assert.NoError(t, err)
		assert.Equal(t, false, isExists)
	})
}

func TestUserRepository_MarkInvitationStatus(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		MailServiceID:  utils.GenerateID(),
		Email:          "ini@email.valid",
		Name:           "ini nama",
		InvitationCode: utils.GenerateRandomAlphanumeric(10),
		Role:           rbac.RoleAdmin,
		Status:         model.InvitationStatusSent,
	}

	t.Run("ok - marked", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_invitations" SET`).WithArgs(invitation.MailServiceID, invitation.Email, invitation.Name, invitation.InvitationCode, invitation.Role, invitation.Status, invitation.ID).WillReturnResult(sqlmock.NewResult(invitation.ID, 1))
		mock.ExpectCommit()

		err := repo.MarkInvitationStatus(ctx, invitation, model.InvitationStatusSent)

		assert.NoError(t, err)
	})

	t.Run("err - from db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "user_invitations" SET`).WithArgs(invitation.MailServiceID, invitation.Email, invitation.Name, invitation.InvitationCode, invitation.Role, invitation.Status, invitation.ID).WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		err := repo.MarkInvitationStatus(ctx, invitation, model.InvitationStatusSent)

		assert.Error(t, err)
	})
}

func TestUserRepository_CheckIsInvitationExists(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		MailServiceID:  utils.GenerateID(),
		Email:          "ini@email.valid",
		Name:           "ini nama",
		InvitationCode: utils.GenerateRandomAlphanumeric(10),
		Role:           rbac.RoleAdmin,
		Status:         model.InvitationStatusSent,
	}

	t.Run("ok - exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "user_invitations" WHERE`).WithArgs(invitation.InvitationCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(invitation.ID))

		err := repo.CheckIsInvitationExists(ctx, invitation.InvitationCode)

		assert.NoError(t, err)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "user_invitations" WHERE`).WithArgs(invitation.InvitationCode).WillReturnError(gorm.ErrRecordNotFound)

		err := repo.CheckIsInvitationExists(ctx, invitation.InvitationCode)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM "user_invitations" WHERE`).WithArgs(invitation.InvitationCode).WillReturnError(errors.New("err db"))

		err := repo.CheckIsInvitationExists(ctx, invitation.InvitationCode)

		assert.Error(t, err)
	})
}
