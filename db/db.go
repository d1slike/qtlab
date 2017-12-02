package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var (
	Db  *gorm.DB
	err error
)

type Doctor struct {
	ID         uint
	Fio        string
	Speciality string
	Cabinet    int
	ClinicID   uint
	Clinic     Clinic
}

type Clinic struct {
	ID      uint
	Name    string
	Address string
	Number  string
	Email   string
}

type Config struct {
	ID     uint
	Width  int
	Height int
}

type FK struct {
	Key uint
}

var Specialities = map[string]interface{}{
	"Врач-колопроктолог":        "Врач-колопроктолог",
	"Врач-офтальмолог":          "Врач-офтальмолог",
	"Врач-пульмонолог":          "Врач-пульмонолог",
	"Врач-ревматолог":           "Врач-ревматолог",
	"Врач-онколог":              "Врач-онколог",
	"Врач-оториноларинголог":    "Врач-оториноларинголог",
	"Врач-хирург":               "Врач-хирург",
	"Врач-эндокринолог":         "Врач-эндокринолог",
	"Врач-терапевт участковый":  "Врач-терапевт участковый",
	"Врач-гастроэнтеролог":      "Врач-гастроэнтеролог",
	"Врач-уролог":               "Врач-уролог",
	"Врач-невролог":             "Врач-невролог",
	"Врач-травматолог-ортопед":  "Врач-травматолог-ортопед",
	"Врач-инфекционист":         "Врач-инфекционист",
	"Врач-кардиолог":            "Врач-кардиолог",
	"Врач-аллерголог-иммунолог": "Врач-аллерголог-иммунолог",
}

func init() {
	Db, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal("Could not connected to sqllite", err)
	}
	Db.AutoMigrate(&Clinic{}, &Doctor{}, &Config{})
}

func GetAllClinics(filter map[string]interface{}) []Clinic {
	var clinics []Clinic
	filterToQuery(filter).Find(&clinics)
	return clinics
}

func GetAllDoctors(filter map[string]interface{}) []Doctor {
	var doctors []Doctor
	filterToQuery(filter).Preload("Clinic").Find(&doctors)
	return doctors
}

func filterToQuery(filter map[string]interface{}) *gorm.DB {
	if filter == nil || len(filter) == 0 {
		return Db
	}
	var query *gorm.DB
	for key, value := range filter {
		query = addWhere(query, key, value)
	}
	return query
}

func addWhere(query *gorm.DB, key string, value interface{}) *gorm.DB {
	if query == nil {
		query = Db
	}
	fk, castFk := value.(FK)
	integer, castInt := value.(int)
	str, castStr := value.(string)
	if castFk {
		query = query.Where(key+" = ?", fk.Key)
	} else if castInt {
		query = query.Where(key+" = ?", integer)
	} else if castStr {
		query = query.Where(key+" LIKE ?", "%"+str+"%")
	}
	return query
}
