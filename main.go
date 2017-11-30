package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"strconv"
	"github.com/therecipe/qt/core"
	"github.com/d1slike/qtlab/form"
	"github.com/jinzhu/gorm"
)

var (
	clinicTable *widgets.QTableWidget
	doctorTable *widgets.QTableWidget

	removeDoctorButton *widgets.QPushButton
	removeClinicButton *widgets.QPushButton

	cancelDoctorFilterButton *widgets.QPushButton
	cancelClinicFilterButton *widgets.QPushButton
)

func main() {
	// Create application
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Qt Lab")
	window.SetMinimumSize2(1000, 800)

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

	refreshClinicTable();
	refreshDoctorsTable()

	// Show the window
	window.Show()

	defer Db.Close()

	// Execute app
	app.Exec()
}

func makeClinicWidget() *widgets.QWidget {
	clinicWidget := widgets.NewQWidget(nil, 0)
	clinicLayout := widgets.NewQVBoxLayout()

	addClinicButton := widgets.NewQPushButton2("Добавить", nil)
	addClinicButton.ConnectClicked(func(c bool) {
		fields := []form.Field{
			{Name: "name", Label: "Название", Type: form.StringType, Required: true},
			{Name: "address", Label: "Адрес", Type: form.StringType},
			{Name: "number", Label: "Контактный номер", Type: form.StringType},
			{Name: "email", Label: "E-mail", Type: form.StringType}}
		form.ShowForm(fields, func(result map[string]interface{}) {
			clinic := Clinic{Name: result["name"].(string)}
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
			Db.Create(&clinic)
			refreshClinicTable()
		}, "Добавить", "Добавление клиники")
	})

	removeClinicButton = widgets.NewQPushButton2("Удалить", nil)
	removeClinicButton.SetEnabled(false)
	removeClinicButton.ConnectClicked(func(checked bool) {
		index := clinicTable.CurrentRow()
		if index >= 0 {
			item := clinicTable.ItemAt2(index, 0)
			id, _ := strconv.Atoi(item.Text())
			Db.Delete(&Clinic{Model: gorm.Model{ID: uint(id)}})
			clinicTable.RemoveRow(index)
			clinicTable.ClearSelection()
		}
	})

	filterClinicButton := widgets.NewQPushButton2("Искать", nil)

	cancelClinicFilterButton = widgets.NewQPushButton2("Отменить фильтр", nil)
	cancelClinicFilterButton.SetEnabled(false)

	clinicButtonsLayout := widgets.NewQHBoxLayout();
	clinicButtonsLayout.AddWidget(addClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(removeClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(filterClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(cancelClinicFilterButton, 0, core.Qt__AlignCenter)

	clinicTable = widgets.NewQTableWidget(nil)
	clinicTable.SetColumnCount(5)
	clinicTable.SetHorizontalHeaderLabels([]string{"ID", "Название", "Адресс", "Номер", "E-mail"})
	clinicTable.SetShowGrid(true)
	clinicTable.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	clinicTable.ConnectSelectRow(func(row int) {
		removeClinicButton.SetEnabled(true)
	})

	clinicLayout.AddWidget(clinicTable, 0, 0)
	clinicLayout.AddLayout(clinicButtonsLayout, 0)
	clinicWidget.SetLayout(clinicLayout)
	return clinicWidget
}

func makeDoctorWidget() *widgets.QWidget {
	doctorWidget := widgets.NewQWidget(nil, 0)
	doctorLayout := widgets.NewQVBoxLayout()

	addDoctorButton := widgets.NewQPushButton2("Добавить", nil)

	removeDoctorButton = widgets.NewQPushButton2("Удалить", nil)
	removeDoctorButton.SetEnabled(false)

	filterDoctorButton := widgets.NewQPushButton2("Искать", nil)

	cancelDoctorFilterButton = widgets.NewQPushButton2("Отменить фильтр", nil)
	cancelDoctorFilterButton.SetEnabled(false)

	doctorButtonsLayout := widgets.NewQHBoxLayout();
	doctorButtonsLayout.AddWidget(addDoctorButton, 0, core.Qt__AlignCenter)
	doctorButtonsLayout.AddWidget(removeDoctorButton, 0, core.Qt__AlignCenter)
	doctorButtonsLayout.AddWidget(filterDoctorButton, 0, core.Qt__AlignCenter)
	doctorButtonsLayout.AddWidget(cancelDoctorFilterButton, 0, core.Qt__AlignCenter)

	doctorTable = widgets.NewQTableWidget(nil)
	doctorTable.SetColumnCount(5)
	doctorTable.SetHorizontalHeaderLabels([]string{"ID", "ФИО", "Специальность", "Кабинет", "Клиника"})
	clinicTable.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	doctorTable.SetShowGrid(true)

	doctorLayout.AddWidget(doctorTable, 0, 0)
	doctorLayout.AddLayout(doctorButtonsLayout, 0)
	doctorWidget.SetLayout(doctorLayout)
	return doctorWidget
}

func refreshClinicTable() {
	clearTable(clinicTable)
	for _, c := range GetAllClinics(nil) {
		clinicTable.InsertRow(0)
		clinicTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(int(c.ID)), 0))
		clinicTable.SetItem(0, 1, widgets.NewQTableWidgetItem2(c.Name, 0))
		clinicTable.SetItem(0, 2, widgets.NewQTableWidgetItem2(c.Address, 0))
		clinicTable.SetItem(0, 3, widgets.NewQTableWidgetItem2(c.Number, 0))
		clinicTable.SetItem(0, 4, widgets.NewQTableWidgetItem2(c.Email, 0))
	}
}

func refreshDoctorsTable() {
	clearTable(doctorTable)
	for _, doctor := range GetAllDoctors(nil) {
		doctorTable.InsertRow(0)
		clinicTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(strconv.Itoa(int(doctor.ID)), 0))
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
