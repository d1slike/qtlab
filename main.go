package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"strconv"
	"github.com/therecipe/qt/core"
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
	window.SetMinimumSize2(800, 800)

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

	removeClinicButton = widgets.NewQPushButton2("Удалить", nil)
	removeClinicButton.SetEnabled(false)

	filterClinicButton := widgets.NewQPushButton2("Искать", nil)

	cancelClinicFilterButton = widgets.NewQPushButton2("Отменить филтр", nil)
	cancelClinicFilterButton.SetEnabled(false)

	clinicButtonsLayout := widgets.NewQHBoxLayout();
	clinicButtonsLayout.AddWidget(addClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(removeClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(filterClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(cancelClinicFilterButton, 0, core.Qt__AlignCenter)

	clinicTable = widgets.NewQTableWidget(nil)
	clinicTable.SetColumnCount(4)
	clinicTable.SetHorizontalHeaderLabels([]string{"Название", "Адресс", "Номер", "E-mail"})
	clinicTable.SetShowGrid(true)

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
	doctorTable.SetColumnCount(4)
	doctorTable.SetHorizontalHeaderLabels([]string{"ФИО", "Специальность", "Кабинет", "Клиника"})
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
		clinicTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(c.Name, 0))
		clinicTable.SetItem(0, 1, widgets.NewQTableWidgetItem2(c.Address, 0))
		clinicTable.SetItem(0, 2, widgets.NewQTableWidgetItem2(c.Number, 0))
		clinicTable.SetItem(0, 3, widgets.NewQTableWidgetItem2(c.Email, 0))
	}
}

func refreshDoctorsTable() {
	clearTable(doctorTable)
	for _, doctor := range GetAllDoctors(nil) {
		doctorTable.InsertRow(0)
		doctorTable.SetItem(0, 0, widgets.NewQTableWidgetItem2(doctor.Fio, 0))
		doctorTable.SetItem(0, 1, widgets.NewQTableWidgetItem2(doctor.Speciality, 0))
		doctorTable.SetItem(0, 2, widgets.NewQTableWidgetItem2(strconv.Itoa(doctor.Cabinet), 0))
		doctorTable.SetItem(0, 3, widgets.NewQTableWidgetItem2(doctor.Clinic.Name, 0))
	}
}

func clearTable(table *widgets.QTableWidget) {
	count := table.RowCount()
	for i := 0; i < count; i++ {
		table.RemoveRow(0)
	}
}
