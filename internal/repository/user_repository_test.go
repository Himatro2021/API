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

func TestUserRepository_GetUserByEmail(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	user := &model.User{
		ID:       utils.GenerateID(),
		Name:     "lucky hehe",
		Email:    "test@mail.he",
		Password: "test password",
		Role:     rbac.RoleAdmin,
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

		res, err := repo.GetUserByEmail(ctx, user.Email)

		assert.NoError(t, err)
		assert.Equal(t, res.ID, user.ID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(user.Email).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetUserByEmail(ctx, user.Email)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(user.Email).
			WillReturnError(errors.New("err db"))

		_, err := repo.GetUserByEmail(ctx, user.Email)

		assert.Error(t, err)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	user := &model.User{
		ID:       utils.GenerateID(),
		Name:     "lucky hehe",
		Email:    "test@mail.he",
		Password: "test password",
		Role:     rbac.RoleAdmin,
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(user.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

		res, err := repo.GetUserByID(ctx, user.ID)

		assert.NoError(t, err)
		assert.Equal(t, res.ID, user.ID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.+ FROM "users"`).
			WithArgs(user.ID).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetUserByID(ctx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.+ FROM "users"`).
			WithArgs(user.ID).
			WillReturnError(errors.New("err db"))

		_, err := repo.GetUserByID(ctx, user.ID)

		assert.Error(t, err)
	})
}

func TestUserRepository_GetUserInvitationByInvitationCode(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	invitation := &model.UserInvitation{
		ID:             utils.GenerateID(),
		MailServiceID:  utils.GenerateID(),
		Email:          "test@email.com",
		Name:           "lucky aja",
		InvitationCode: "sembaranganKodeAja",
		Role:           rbac.RoleAdmin,
		Status:         model.InvitationStatusSent,
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "user_invitations"`).
			WithArgs(invitation.InvitationCode).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(invitation.ID))

		res, err := repo.GetUserInvitationByInvitationCode(ctx, invitation.InvitationCode)

		assert.NoError(t, err)
		assert.Equal(t, res.ID, invitation.ID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "user_invitations"`).
			WithArgs(invitation.InvitationCode).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetUserInvitationByInvitationCode(ctx, invitation.InvitationCode)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "user_invitations"`).
			WithArgs(invitation.InvitationCode).
			WillReturnError(errors.New("err db"))

		_, err := repo.GetUserInvitationByInvitationCode(ctx, invitation.InvitationCode)

		assert.Error(t, err)
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

func TestUserRepository_Register(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	repo := userRepository{
		db: kit.db,
	}

	user := &model.User{
		ID:       utils.GenerateID(),
		Name:     "lucky hehe",
		Email:    "test@mail.he",
		Password: "test password",
		Role:     rbac.RoleAdmin,
	}

	t.Run("ok - registered", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "users"`).
			WithArgs(user.Name, user.Email, user.Password, user.Role, user.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
		mock.ExpectCommit()

		err := repo.Register(ctx, user)

		assert.NoError(t, err)
	})

	t.Run("err - err db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "users"`).
			WithArgs(user.Name, user.Email, user.Password, user.Role, user.ID).
			WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		err := repo.Register(ctx, user)

		assert.Error(t, err)
	})
}
