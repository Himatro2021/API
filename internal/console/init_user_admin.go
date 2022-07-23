package console

import (
	"github.com/Himatro2021/API/internal/db"
	"github.com/Himatro2021/API/internal/helper"
	"github.com/Himatro2021/API/internal/model"
	"github.com/Himatro2021/API/internal/rbac"
	"github.com/go-playground/validator/v10"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var initAdminCmd = &cobra.Command{
	Use:     "init-admin",
	Short:   "initialize admin user",
	Long:    "initialize credentials to create an user with admin role",
	Example: "go run main.go init-admin your@email.com yourVerySecretPasswordDontForgetThis!",
	Run:     initAdmin,
}

func init() {
	RootCmd.AddCommand(initAdminCmd)
}

func initAdmin(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		logrus.Fatal("invalid command usage")
	}

	type input struct {
		Email    string `validate:"email"`
		Password string `validate:"required,min=8"`
	}

	data := input{
		Email:    args[0],
		Password: args[1],
	}

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		logrus.Fatal("input error:", err.Error())
	}

	db.InitializePostgresConn()

	fatalIfAdminAlreadyPresent()

	hashed, err := helper.HashString(data.Password)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// TODO hash
	user := &model.User{
		ID:       utils.GenerateID(),
		Name:     "Super Admin",
		Email:    data.Email,
		Password: hashed,
	}
	user.SetRole(rbac.RoleAdmin)

	err = db.PostgresDB.Model(&model.User{}).Create(user).Error
	if err != nil {
		logrus.Fatal("unexpected error:", err.Error())
	}
}

func fatalIfAdminAlreadyPresent() {
	user := &model.User{}

	err := db.PostgresDB.Model(&model.User{}).Where("role = 'ADMIN'::user_roles").Take(user).Error
	switch err {
	default:
		logrus.Fatal("unexpected error:", err.Error())
	case nil:
		logrus.Fatal("User as admin already exits")
	case gorm.ErrRecordNotFound:
		return
	}
}
