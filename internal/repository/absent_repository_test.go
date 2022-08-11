package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Himatro2021/API/internal/db"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/kumparan/go-utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAbsentRepository_GetAbsentFormByID(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	formID := utils.GenerateID()
	repo := absentRepository{
		db: kit.db,
	}

	t.Run("ok - found as member", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).
			WithArgs(formID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(formID))

		absentForm, err := repo.GetAbsentFormByID(ctx, formID)

		assert.NoError(t, err)
		assert.Equal(t, absentForm.ID, formID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).
			WithArgs(formID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, err := repo.GetAbsentFormByID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).
			WithArgs(formID).WillReturnError(errors.New("err from db"))

		_, err := repo.GetAbsentFormByID(ctx, formID)

		assert.Error(t, err)
	})
}

func TestAbsentRepository_GetAbsentListByID(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	formID := utils.GenerateID()
	listID := utils.GenerateID()
	repo := absentRepository{
		db: kit.db,
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(listID, formID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "absent_form_id"}).AddRow(listID, formID))

		absentList, err := repo.GetAbsentListByID(ctx, formID, listID)

		assert.NoError(t, err)
		assert.Equal(t, absentList.ID, listID)
		assert.Equal(t, absentList.AbsentFormID, formID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(listID, formID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "absent_form_id"}))

		_, err := repo.GetAbsentListByID(ctx, formID, listID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(listID, formID).WillReturnError(errors.New("err db"))

		_, err := repo.GetAbsentListByID(ctx, formID, listID)

		assert.Error(t, err)
	})
}

func TestAbsentRepo_GetAbsentListByCreatorID(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	LoadConf()

	ctx := context.TODO()
	creatorID := utils.GenerateID()
	formID := utils.GenerateID()
	repo := absentRepository{
		db: kit.db,
	}

	t.Run("ok - found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(formID, creatorID).
			WillReturnRows(sqlmock.NewRows([]string{"absent_form_id", "created_by"}).AddRow(formID, creatorID))

		absentList, err := repo.GetAbsentListByCreatorID(ctx, formID, creatorID)

		assert.NoError(t, err)
		assert.NotEqual(t, absentList, nil)
		assert.Equal(t, absentList.AbsentFormID, formID)
		assert.Equal(t, absentList.CreatedBy, creatorID)
	})

	t.Run("ok - not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(formID, creatorID).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetAbsentListByCreatorID(ctx, formID, creatorID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_lists"`).
			WithArgs(formID, creatorID).
			WillReturnError(errors.New("err from db"))

		_, err := repo.GetAbsentListByCreatorID(ctx, formID, creatorID)

		assert.Error(t, err)
	})
}

func TestAbsentRepo_CreateAbsentForm(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock
	ctx := context.TODO()
	repo := absentRepository{
		db: kit.db,
	}

	title := "ini title"
	start, _ := helper.ParseDateAndTimeStringToTime("2001-10-29", "12:00")
	finish, _ := helper.ParseDateAndTimeStringToTime("3001-1-20", "13:00")
	groupID := utils.GenerateID()
	userID := utils.GenerateID()
	form := &model.AbsentForm{
		ID:         utils.GenerateID(),
		Title:      title,
		StartAt:    start,
		FinishedAt: finish,
	}

	t.Run("ok - created", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "absent_forms"`).WithArgs(groupID, start, finish, title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(form.ID))
		mock.ExpectCommit()

		absentForm, err := repo.CreateAbsentForm(ctx, title, start, finish, groupID, userID)

		assert.NoError(t, err)
		assert.Equal(t, absentForm.ParticipantGroupID, groupID)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "absent_forms"`).WithArgs(groupID, start, finish, title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		_, err := repo.CreateAbsentForm(ctx, title, start, finish, groupID, userID)

		assert.Error(t, err)
	})
}

func TestAbsentRepository_FillAbsentFormByAttendee(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock
	ctx := context.TODO()
	repo := absentRepository{
		db: kit.db,
	}

	userID := utils.GenerateID()
	formID := utils.GenerateID()
	status := model.Present
	reason := "empty"

	t.Run("ok - filled", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "absent_lists"`).WithArgs(formID, userID, sqlmock.AnyArg(), status, sqlmock.AnyArg(), reason, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"absent_form_id"}).AddRow(formID))
		mock.ExpectCommit()

		absentList, err := repo.FillAbsentFormByAttendee(ctx, userID, formID, status, reason)

		assert.NoError(t, err)
		assert.Equal(t, absentList.AbsentFormID, formID)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "absent_lists"`).WithArgs(formID, userID, sqlmock.AnyArg(), status, sqlmock.AnyArg(), reason, sqlmock.AnyArg()).
			WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		_, err := repo.FillAbsentFormByAttendee(ctx, userID, formID, status, reason)

		assert.Error(t, err)
	})
}

func TestAbsentRepository_UpdateAbsentForm(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock
	ctx := context.TODO()
	repo := absentRepository{
		db: kit.db,
	}

	title := "ini title"
	start, _ := helper.ParseDateAndTimeStringToTime("2001-10-29", "12:00")
	finish, _ := helper.ParseDateAndTimeStringToTime("3001-1-20", "13:00")
	groupID := utils.GenerateID()
	userID := utils.GenerateID()
	form := &model.AbsentForm{
		ID:         utils.GenerateID(),
		Title:      title,
		StartAt:    start,
		FinishedAt: finish,
	}

	t.Run("ok - created", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).WithArgs(form.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(form.ID))

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "absent_forms" SET`).
			WithArgs(groupID, start, finish, title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(form.ID, 1))
		mock.ExpectCommit()

		absentForm, err := repo.UpdateAbsentForm(ctx, form.ID, title, start, finish, groupID, userID)

		assert.NoError(t, err)
		assert.Equal(t, absentForm.ID, form.ID)
	})

	t.Run("ok - absent form not found", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).WithArgs(form.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, err := repo.UpdateAbsentForm(ctx, form.ID, title, start, finish, groupID, userID)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err when find absent form", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).WithArgs(form.ID).WillReturnError(errors.New("err db"))

		_, err := repo.UpdateAbsentForm(ctx, form.ID, title, start, finish, groupID, userID)

		assert.Error(t, err)
	})

	t.Run("err when updating", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT .+ FROM "absent_forms"`).WithArgs(form.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(form.ID))

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "absent_forms" SET`).
			WithArgs(groupID, start, finish, title, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("err db"))
		mock.ExpectCommit()

		_, err := repo.UpdateAbsentForm(ctx, form.ID, title, start, finish, groupID, userID)

		assert.Error(t, err)
	})
}

func TestAbsentRepository_UpdateAbsentListByAttendee(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock
	ctx := context.TODO()
	repo := absentRepository{
		db: kit.db,
	}

	absentList := &model.AbsentList{
		ID:            utils.GenerateID(),
		AbsentFormID:  utils.GenerateID(),
		CreatedBy:     int64(1),
		Status:        model.Absent,
		ExecuseReason: "",
	}

	t.Run("ok - updated", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "absent_lists" SET`).WithArgs(absentList.AbsentFormID, absentList.CreatedBy, sqlmock.AnyArg(), sqlmock.AnyArg(), absentList.Status, absentList.ExecuseReason, absentList.ID).WillReturnResult(sqlmock.NewResult(absentList.ID, 1))
		mock.ExpectCommit()

		result, err := repo.UpdateAbsentListByAttendee(ctx, absentList)

		assert.NoError(t, err)
		assert.Equal(t, result.ID, absentList.ID)
	})

	t.Run("err when updating", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "absent_lists" SET`).WithArgs(absentList.AbsentFormID, absentList.CreatedBy, sqlmock.AnyArg(), sqlmock.AnyArg(), absentList.Status, absentList.ExecuseReason, absentList.ID).WillReturnError(errors.New("err db"))
		mock.ExpectRollback()

		_, err := repo.UpdateAbsentListByAttendee(ctx, absentList)

		assert.Error(t, err)
	})
}

func TestAbsentRepository_GetParticipantsByFormID(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock
	ctx := context.TODO()
	repo := absentRepository{
		db: kit.db,
	}

	formID := utils.GenerateID()
	_ = os.Setenv("PRIVATE_KEY", "supersecret")
	_ = os.Setenv("IV_KEY", "4e6064d3814c2cd22c550155655fefc6")

	t.Run("ok", func(t *testing.T) {
		mock.ExpectQuery(`^select .+ from absent_lists al inner join users u ON u.id = al.created_by where al.absent_form_id`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("32dd72fd391d8e334d6de3e5fce4cfcc"))

		result, err := repo.GetParticipantsByFormID(ctx, formID)

		assert.NoError(t, err)
		assert.NotEqual(t, []model.Participant{}, result)
	})

	t.Run("err from db", func(t *testing.T) {
		mock.ExpectQuery(`^select .+ from absent_lists al inner join users u ON u.id = al.created_by where al.absent_form_id`).
			WillReturnError(errors.New("err db"))

		result, err := repo.GetParticipantsByFormID(ctx, formID)

		assert.Error(t, err)
		assert.Equal(t, []model.Participant{}, result)
	})
}

func TestAbsentRepository_GetAbsentResultFromCache(t *testing.T) {
	ctx := context.TODO()
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cacher := db.NewCacher(client)
	repo := absentRepository{
		cacher: cacher,
	}

	absentResult := &model.AbsentResult{
		Title:      "test doang",
		StartAt:    time.Now().Add(-1 * time.Hour),
		FinishedAt: time.Now().Add(1 * time.Hour),
		Participants: []model.Participant{
			{
				Name:   "lucky test",
				Status: model.Present,
				Reason: "",
			},
		},
	}
	absentResultByte, err := json.Marshal(absentResult)
	cacheKey := model.AbsentResult.CacheKeyByFormID(model.AbsentResult{}, utils.GenerateID())
	assert.NoError(t, err)

	t.Run("ok - found", func(t *testing.T) {
		err := cacher.Set(ctx, cacheKey, string(absentResultByte), time.Second*100)
		assert.NoError(t, err)

		response, err := repo.GetAbsentResultFromCache(ctx, cacheKey)
		assert.NoError(t, err)
		assert.Equal(t, absentResult.Title, response.Title)

		mr.FlushDB()
	})

	t.Run("ok - cache not found", func(t *testing.T) {
		_, err := repo.GetAbsentResultFromCache(ctx, cacheKey)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	// close redis to simulate redis error
	mr.Close()

	t.Run("err - redis error", func(t *testing.T) {
		_, err := repo.GetAbsentResultFromCache(ctx, cacheKey)

		assert.Error(t, err)
	})
}

func TestAbsentRepo_SetAbsentResultToCache(t *testing.T) {
	ctx := context.TODO()
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cacher := db.NewCacher(client)
	repo := absentRepository{
		cacher: cacher,
	}

	absentResult := &model.AbsentResult{
		Title:      "test doang",
		StartAt:    time.Now().Add(-1 * time.Hour),
		FinishedAt: time.Now().Add(1 * time.Hour),
		Participants: []model.Participant{
			{
				Name:   "lucky test",
				Status: model.Present,
				Reason: "",
			},
		},
	}
	absentResultByte, err := json.Marshal(absentResult)
	assert.NoError(t, err)
	formID := utils.GenerateID()
	cacheKey := model.AbsentResult.CacheKeyByFormID(model.AbsentResult{}, formID)

	t.Run("ok - set to cache", func(t *testing.T) {
		err := repo.SetAbsentResultToCache(ctx, absentResult, formID)
		assert.NoError(t, err)

		resp, err := mr.Get(cacheKey)

		assert.NoError(t, err)
		assert.Equal(t, resp, string(absentResultByte))

		mr.FlushDB()
	})

	t.Run("redis error", func(t *testing.T) {
		mr.Close()

		err := repo.SetAbsentResultToCache(ctx, absentResult, formID)

		assert.Error(t, err)
	})
}

func TestAbsentRepo_UpdateParticipantsInAbsentResultCache(t *testing.T) {
	kit := initializeRepoTestKit(t)
	mock := kit.dbmock

	ctx := context.TODO()
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	_ = os.Setenv("PRIVATE_KEY", "supersecret")
	_ = os.Setenv("IV_KEY", "4e6064d3814c2cd22c550155655fefc6")

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cacher := db.NewCacher(client)
	repo := absentRepository{
		cacher: cacher,
		db:     kit.db,
	}

	absentResult := &model.AbsentResult{
		Title:      "test doang",
		StartAt:    time.Now().Add(-1 * time.Hour),
		FinishedAt: time.Now().Add(1 * time.Hour),
		Participants: []model.Participant{
			{
				Name:   "lucky test",
				Status: model.Present,
				Reason: "",
			},
		},
	}

	absentResultByte, err := json.Marshal(absentResult)
	assert.NoError(t, err)
	formID := utils.GenerateID()
	cacheKey := model.AbsentResult.CacheKeyByFormID(model.AbsentResult{}, formID)

	t.Run("ok - updated", func(t *testing.T) {
		err := mr.Set(cacheKey, string(absentResultByte))
		assert.NoError(t, err)
		defer mr.FlushDB()

		mock.ExpectQuery(`^select .+ from absent_lists al inner join users u ON u.id = al.created_by where al.absent_form_id`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("32dd72fd391d8e334d6de3e5fce4cfcc").AddRow("32dd72fd391d8e334d6de3e5fce4cfcc"))

		err = repo.UpdateParticipantsInAbsentResultCache(ctx, cacheKey)
		assert.NoError(t, err)

		resp, err := mr.Get(cacheKey)
		assert.NoError(t, err)

		result := &model.AbsentResult{}
		err = json.Unmarshal([]byte(resp), result)
		assert.NoError(t, err)

		assert.Equal(t, len(result.Participants), 2)
	})

	t.Run("ok - cache not found", func(t *testing.T) {
		err := repo.UpdateParticipantsInAbsentResultCache(ctx, cacheKey)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})

	t.Run("err - redis err", func(t *testing.T) {
		mr.Close()

		err := repo.UpdateParticipantsInAbsentResultCache(ctx, cacheKey)

		assert.Error(t, err)
	})

	err = mr.Restart()
	assert.NoError(t, err)

	t.Run("err - db err when getting participants", func(t *testing.T) {
		err := mr.Set(cacheKey, string(absentResultByte))
		assert.NoError(t, err)
		defer mr.FlushDB()

		mock.ExpectQuery(`^select .+ from absent_lists al inner join users u ON u.id = al.created_by where al.absent_form_id`).
			WillReturnError(errors.New("err db"))

		err = repo.UpdateParticipantsInAbsentResultCache(ctx, cacheKey)

		assert.Error(t, err)
	})

	// TODO how to simulate error when setting new participants info in the cache?
}
