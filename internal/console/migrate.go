package console

import (
	"himatro-api/internal/db"
	"himatro-api/internal/models"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate all database table",
	Long:  "Use this command to initialize your database table for the first time",
	Run:   migrate,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	db.Connect()

	db.DB.AutoMigrate(&models.AnggotaBiasa{})
	db.DB.AutoMigrate(&models.Jabatan{})
	db.DB.AutoMigrate(&models.Pengurus{})
	db.DB.AutoMigrate(&models.FormAbsensi{})
	db.DB.AutoMigrate(&models.AbsentList{})
	db.DB.AutoMigrate(&models.Departemen{})
	db.DB.AutoMigrate(&models.User{})
}
