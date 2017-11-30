package ui

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/core"
	"strconv"
)

func MakeWidget(
	removeButton *widgets.QPushButton,
	cancelFilterButton *widgets.QPushButton,
	table *widgets.QTableWidget,
	header []string,
	onClinicAdd func(bool),
	onDelete func(id uint),
	onClickFilter func(bool),
	onClickCancelFilter func()) *widgets.QWidget {

	clinicWidget := widgets.NewQWidget(nil, 0)
	clinicLayout := widgets.NewQVBoxLayout()

	addClinicButton := widgets.NewQPushButton2("Добавить", nil)
	addClinicButton.ConnectClicked(onClinicAdd)

	removeButton.SetEnabled(false)
	removeButton.ConnectClicked(func(checked bool) {
		index := table.CurrentRow()
		if index >= 0 {
			item := table.ItemAt2(index, 0)
			id, _ := strconv.Atoi(item.Text())
			onDelete(uint(id))
			table.RemoveRow(index)
			table.ClearSelection()
		}
	})

	filterClinicButton := widgets.NewQPushButton2("Искать", nil)
	filterClinicButton.ConnectClicked(onClickFilter)

	cancelFilterButton.SetEnabled(false)
	cancelFilterButton.ConnectClicked(func(checked bool) {
		onClickCancelFilter()
		cancelFilterButton.SetEnabled(false)
	})

	clinicButtonsLayout := widgets.NewQHBoxLayout();
	clinicButtonsLayout.AddWidget(addClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(removeButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(filterClinicButton, 0, core.Qt__AlignCenter)
	clinicButtonsLayout.AddWidget(cancelFilterButton, 0, core.Qt__AlignCenter)

	table.SetColumnCount(len(header))
	table.SetHorizontalHeaderLabels(header)
	table.SetShowGrid(true)
	table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	table.ConnectSelectRow(func(row int) {
		removeButton.SetEnabled(true)
	})

	clinicLayout.AddWidget(table, 0, 0)
	clinicLayout.AddLayout(clinicButtonsLayout, 0)
	clinicWidget.SetLayout(clinicLayout)
	return clinicWidget
}
