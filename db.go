package main

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
	gorm.Model
	Fio        string
	Speciality string
	Cabinet    int
	Clinic     Clinic
}

type Clinic struct {
	gorm.Model
	Name    string
	Address string
	Number  string
	Email   string
}

func init() {
	Db, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal("Could not connected to sqllite", err)
	}
	Db.AutoMigrate(&Clinic{}, &Doctor{})
}

func GetAllClinics(filter map[string]interface{}) []Clinic {
	var clinics []Clinic
	if filter != nil {
		Db.Where(filter).Find(&clinics)
	} else {
		Db.Find(&clinics)
	}
	return clinics
}

func GetAllDoctors(filter map[string]interface{}) []Doctor {
	var doctors []Doctor
	if filter != nil {
		Db.Where(filter).Find(&doctors)
	} else {
		Db.Find(&doctors)
	}
	return doctors
}
