package console

import (
	"encoding/csv"
	"fmt"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var seederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "Run the seeder program",
	Long:  "Use this command to seed your database",
	Run:   seeder,
}

func init() {
	RootCmd.AddCommand(seederCmd)
}

func seeder(cmd *cobra.Command, args []string) {
	db.Connect()
	seedAnggotaBiasa()
	seedDepartemen()
	seedJabatan()
	seedPengurus()
	seedSuperAdminUser()
}

func seedAnggotaBiasa() {
	filepath := os.Getenv("ANGGOTA_BIASA_SEEDER_DATA_PATH")
	data := readFromCSV(filepath)

	validateCSVColumnHeader(data[0], "ANGGOTA_BIASA_CSV_HEADER_CONFIG")

	for i := 1; i < len(data); i++ {
		anggotaBiasa := models.AnggotaBiasa{
			NPM:  data[i][0],
			Nama: data[i][1],
		}

		log.Printf("Inserting Anggota Biasa data with NPM: %s, Nama: %s", anggotaBiasa.NPM, anggotaBiasa.Nama)

		result := db.DB.Create(&anggotaBiasa)

		if result.Error != nil {
			log.Printf("Failed to insert data for %q because: %s", anggotaBiasa, result.Error)
		}
	}
}

func seedDepartemen() {
	filepath := os.Getenv("DEPARTEMEN_SEEDER_DATA_PATH")
	data := readFromCSV(filepath)

	validateCSVColumnHeader(data[0], "DEPARTEMEN_CSV_HEADER_CONFIG")

	for i := 1; i < len(data); i++ {
		id, err := strconv.ParseInt(data[i][0], 10, 8)

		if err != nil {
			panic("Invalid value in departemenID from corespondense CSV file.")
		}

		departemen := models.Departemen{
			ID:   int(id),
			Nama: data[i][1],
		}

		log.Printf("Inserting Departemen data with ID: %d, Nama: %s", departemen.ID, departemen.Nama)

		result := db.DB.Create(&departemen)

		if result.Error != nil {
			log.Printf("Failed to insert data for %q because: %s", departemen, result.Error)
		}
	}
}

func seedPengurus() {
	filepath := os.Getenv("PENGURUS_SEEDER_DATA_PATH")
	data := readFromCSV(filepath)

	validateCSVColumnHeader(data[0], "PENGURUS_CSV_HEADER_CONFIG")

	for i := 1; i < len(data); i++ {
		departemenID, err := strconv.Atoi(data[i][1])

		if err != nil {
			panic("Invalid departemenID data in pengurus.csv file.")
		}

		jabatanID, err := strconv.Atoi(data[i][2])

		if err != nil {
			panic("Invalid jabatanID data in pengurus.csv file.")
		}

		pengurus := models.Pengurus{
			NPM:          data[i][0],
			DepartemenID: departemenID,
			JabatanID:    jabatanID,
		}

		log.Printf("Inserting data Pengurus with NPM: %s", pengurus.NPM)

		result := db.DB.Create(&pengurus)

		if result.Error != nil {
			log.Printf("Failed to insert data Pengurus for %s because %s", pengurus.NPM, result.Error)
		}
	}
}

func seedJabatan() {
	filepath := os.Getenv("JABATAN_SEEDER_DATA_PATH")
	data := readFromCSV(filepath)

	validateCSVColumnHeader(data[0], "JABATAN_CSV_HEADER_CONFIG")

	for i := 1; i < len(data); i++ {
		id, err := strconv.ParseInt(data[i][0], 10, 8)

		if err != nil {
			panic("Invalid type on JabatanID field from corespondense CSV file")
		}

		privLevel, err := strconv.ParseInt(data[i][1], 10, 8)

		if err != nil {
			panic("Invalid type on privilegeLevel field from corespondense CSV file")
		}

		jabatan := models.Jabatan{
			ID:             uint(id),
			PrivilegeLevel: int(privLevel),
			Name:           data[i][2],
		}

		result := db.DB.Create(&jabatan)

		if result.Error != nil {
			log.Printf("Failed to insert data Jabatan for %s", jabatan.Name)
		}
	}
}

func seedSuperAdminUser() {
	filepath := os.Getenv("SUPER_ADMIN_SEEDER_DATA_PATH")
	data := readFromCSV(filepath)

	validateCSVColumnHeader(data[0], "SUPER_ADMIN_CSV_HEADER_CONFIG")

	for i := 1; i < len(data); i++ {
		superAdmin := models.User{
			NPM:      data[i][0],
			Password: data[i][1],
		}

		result := db.DB.Create(&superAdmin)

		if result.Error != nil {
			log.Printf("User admin failed to insert in: %s", superAdmin.NPM)
		}
	}
}

func readFromCSV(filepath string) [][]string {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal("Failed to open CSV seeder file: ", filepath)
	}

	defer file.Close()

	csvData := csv.NewReader(file)

	data, err := csvData.ReadAll()

	if err != nil {
		log.Fatal("error accoured when reading CSV seeder file: ", filepath)
	}

	return data
}

func validateCSVColumnHeader(firstRow []string, configName string) {
	config := strings.Split(os.Getenv(configName), ",")

	if len(config) != len(firstRow) {
		log.Fatal(fmt.Sprintf("Format mismatch in %s with the input CS file. Please read the seeder instruction carefully.", configName))
	}

	for i, s := range config {
		if s != firstRow[i] {
			log.Fatal(fmt.Sprintf("CSV header format mismatch: %s with %s while checking validity in: %s", s, firstRow[i], configName))
		}
	}
}
