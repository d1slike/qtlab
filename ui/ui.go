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

	widget := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout()

	addButton := widgets.NewQPushButton2("Добавить", nil)
	addButton.ConnectClicked(onClinicAdd)

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

	filterButton := widgets.NewQPushButton2("Искать", nil)
	filterButton.ConnectClicked(onClickFilter)

	cancelFilterButton.SetEnabled(false)
	cancelFilterButton.ConnectClicked(func(checked bool) {
		onClickCancelFilter()
		cancelFilterButton.SetEnabled(false)
	})

	buttonsLayout := widgets.NewQHBoxLayout();
	buttonsLayout.AddWidget(addButton, 0, core.Qt__AlignCenter)
	buttonsLayout.AddWidget(removeButton, 0, core.Qt__AlignCenter)
	buttonsLayout.AddWidget(filterButton, 0, core.Qt__AlignCenter)
	buttonsLayout.AddWidget(cancelFilterButton, 0, core.Qt__AlignCenter)

	table.SetColumnCount(len(header))
	table.SetHorizontalHeaderLabels(header)
	table.SetShowGrid(true)
	table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	table.ConnectSelectRow(func(row int) {
		removeButton.SetEnabled(true)
	})

	layout.AddWidget(table, 0, 0)
	layout.AddLayout(buttonsLayout, 0)
	widget.SetLayout(layout)
	return widget
}
