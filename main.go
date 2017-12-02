package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"strconv"
	"github.com/d1slike/qtlab/form"
	"github.com/d1slike/qtlab/db"
	"github.com/d1slike/qtlab/ui"
)

var (
	clinicTable *widgets.QTableWidget
	doctorTable *widgets.QTableWidget

	removeDoctorButton *widgets.QPushButton
	removeClinicButton *widgets.QPushButton

	cancelDoctorFilterButton *widgets.QPushButton
	cancelClinicFilterButton *widgets.QPushButton

	window *widgets.QMainWindow
)

func main() {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create main window
	window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Qt Lab")

	cfg := db.Config{}
	db.Db.First(&cfg)
	if cfg.ID == 0 {
		cfg.Width = 1000
		cfg.Height = 800
	}
	window.Resize2(cfg.Width, cfg.Height)

	// Create main layout
	layout := widgets.NewQVBoxLayout()

	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetLayout(layout)

	tabWidget := widgets.NewQTabWidget(nil)
	tabWidget.AddTab(makeClinicWidget(), "Клиники")
	tabWidget.AddTab(makeDoctorWidget(), "Доктора")

	layout.AddWidget(tabWidget, 0, 0)

	// Set main widget as the central widget of the window
	window.SetCentralWidget(mainWidget)

	refreshClinicTable(db.GetAllClinics(nil))
	refreshDoctorsTable(db.GetAllDoctors(nil))

	// Show the window
	window.Show()

	defer db.Db.Close()

	// Execute app
	app.Exec()

	newCfg := db.Config{}
	db.Db.First(&newCfg)
	newCfg.Width = window.Width()
	newCfg.Height = window.Height()
	db.Db.Save(&newCfg)
}

func makeClinicWidget() *widgets.QWidget {
	clinicTable = widgets.NewQTableWidget(nil)
	cancelClinicFilterButton = widgets.NewQPushButton2("Отменить фильтр", nil)
	removeClinicButton = widgets.NewQPushButton2("Удалить", nil)
	return ui.MakeWidget(
		removeClinicButton,
		cancelClinicFilterButton,
		clinicTable,
		[]string{"ID", "Название", "Адресс", "Номер", "E-mail"},
		func(c bool) {
			fields := []form.Field{
				{Name: "name", Label: "Название", Type: form.StringType, Required: true},
				{Name: "address", Label: "Адрес", Type: form.StringType},
				{Name: "number", Label: "Контактный номер", Type: form.StringType},
				{Name: "email", Label: "E-mail", Type: form.StringType}}
			form.ShowForm(fields, func(result map[string]interface{}) {
				clinic := db.Clinic{Name: result["name"].(string)}
				address, hasAddress := result["address"]
				if hasAddress {
					clinic.Address = address.(string)
				}
				number, hasNumber := result["number"]
				if hasNumber {
					clinic.Number = number.(string)
				}
				email, hasEmail := result["email"]
				if hasEmail {
					clinic.Email = email.(string)
				}
				db.Db.Create(&clinic)
				refreshClinicTable(db.GetAllClinics(nil))
			}, "Добавить", "Добавление клиники", window)
		}, func(id uint) {
			db.Db.Delete(&db.Clinic{ID: id})
		}, func(checked bool) {
			fields := []form.Field{
				{Name: "name", Type: form.StringType, Label: "Название"},
				{Name: "address", Type: form.StringType, Label: "Адрес"},
				{Name: "number", Type: form.StringType, Label: "Телефон"},
				{Name: "email", Type: form.StringType, Label: "E-mail"}}
			form.ShowForm(fields, func(filter map[string]interface{}) {
				refreshClinicTable(db.GetAllClinics(filter))
				cancelClinicFilterButton.SetEnabled(true)
			}, "Искать", "Поиск клиник", window)
		}, func() {
			refreshClinicTable(db.GetAllClinics(nil))
		})
}

func makeDoctorWidget() *widgets.QWidget {
	doctorTable = widgets.NewQTableWidget(nil)
	cancelDoctorFilterButton = widgets.NewQPushButton2("Отменить фильтр", nil)
	removeDoctorButton = widgets.NewQPushButton2("Удалить", nil)
	return ui.MakeWidget(
		removeDoctorButton,
		cancelDoctorFilterButton,
		doctorTable,
		[]string{"ID", "ФИО", "Специальность", "Кабинет", "Клиника"},
		func(c bool) {
			clinicOptions := make(map[string]interface{})
			for _, clinic := range db.GetAllClinics(nil) {
				clinicOptions[clinic.Name] = clinic.ID
			}
			fields := []form.Field{
				{Name: "fio", Label: "ФИО", Type: form.StringType, Required: true},
				{Name: "speciality", Label: "Специальность", Type: form.ObjectType, Required: true, Options: db.Specialities},
				{Name: "cabinet", Label: "Кабинет", Type: form.IntegerType, Required: true},
				{Name: "clinic", Label: "Клиника", Type: form.ObjectType, Required: true, Options: clinicOptions}}
			form.ShowForm(fields, func(result map[string]interface{}) {
				doctor := db.Doctor{Fio: result["fio"].(string),
					Speciality: result["speciality"].(string),
					Cabinet: result["cabinet"].(int),
					ClinicID: result["clinic"].(uint)}
				db.Db.Create(&doctor)
				refreshDoctorsTable(db.GetAllDoctors(nil))
			}, "Добавить", "Добавление врача", window)
		}, func(id uint) {
			db.Db.Delete(&db.Doctor{ID: id})
		}, func(checked bool) {
			clinicOptions := make(map[string]interface{})
			for _, clinic := range db.GetAllClinics(nil) {
				clinicOptions[clinic.Name] = db.FK{Key: clinic.ID}
			}
			fields := []form.Field{
				{Name: "fio", Type: form.StringType, Label: "ФИО",},
				{Name: "speciality", Type: form.ObjectType, Label: "Специальность", Options: db.Specialities},
				{Name: "cabinet", Type: form.IntegerType, Label: "Кабинет"},
				{Name: "clinic_id", Type: form.ObjectType, Label: "Клиника", Options: clinicOptions}}
			form.ShowForm(fields, func(filter map[string]interface{}) {
				refreshDoctorsTable(db.GetAllDoctors(filter))
				cancelDoctorFilterButton.SetEnabled(true)
			}, "Искать", "Поиск врачей", window)
		}, func() {
			refreshDoctorsTable(db.GetAllDoctors(nil))
		})
}

func refreshClinicTable(clinics []db.Clinic) {
	clearTable(clinicTable)
	for _, c := range clinics {
		clinicTable.InsertRow(0)
		clinicTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(int(c.ID)), 0))
		clinicTable.SetItem(0, 1, widgets.NewQTableWidgetItem2(c.Name, 0))
		clinicTable.SetItem(0, 2, widgets.NewQTableWidgetItem2(c.Address, 0))
		clinicTable.SetItem(0, 3, widgets.NewQTableWidgetItem2(c.Number, 0))
		clinicTable.SetItem(0, 4, widgets.NewQTableWidgetItem2(c.Email, 0))
	}
}

func refreshDoctorsTable(doctors []db.Doctor) {
	clearTable(doctorTable)
	for _, doctor := range doctors {
		doctorTable.InsertRow(0)
		doctorTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(int(doctor.ID)), 0))
		doctorTable.SetItem(0, 1, widgets.NewQTableWidgetItem2(doctor.Fio, 0))
		doctorTable.SetItem(0, 2, widgets.NewQTableWidgetItem2(doctor.Speciality, 0))
		doctorTable.SetItem(0, 3, widgets.NewQTableWidgetItem2(strconv.Itoa(doctor.Cabinet), 0))
		doctorTable.SetItem(0, 4, widgets.NewQTableWidgetItem2(doctor.Clinic.Name, 0))
	}
}

func clearTable(table *widgets.QTableWidget) {
	count := table.RowCount()
	for i := 0; i < count; i++ {
		table.RemoveRow(0)
	}
}
