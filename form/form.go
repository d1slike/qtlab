package form

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
	"strings"
	"strconv"
)

const (
	StringType  = "string"
	IntegerType = "integer"
	ObjectType  = "object"
)

type Field struct {
	Name     string
	Label    string
	Type     string
	Required bool
	Options  map[string]interface{}
}

func ShowForm(fields []Field, onClose func(map[string]interface{}), buttonText string, title string, parent widgets.QWidget_ITF) {
	dialog := widgets.NewQDialog(parent, core.Qt__Dialog)
	layout := widgets.NewQVBoxLayout()
	dialog.SetLayout(layout)
	dialog.SetWindowTitle(title)
	dialog.SetMinimumWidth(400)
	result := make(map[string]interface{})
	dialog.ConnectAccepted(func() {
		if len(result) > 0 {
			onClose(result)
		}
	})
	for _, f := range fields {
		var field widgets.QWidget_ITF
		switch f.Type {
		case ObjectType:
			comboBox := widgets.NewQComboBox(nil)
			if f.Options != nil {
				keys := make([]string, len(f.Options)+1)
				keys[0] = " "
				i := 1
				for key := range f.Options {
					keys[i] = key
					i++
				}
				comboBox.AddItems(keys)
				comboBox.SetCurrentIndex(0)
			}
			comboBox.ConnectCurrentIndexChanged2(onChangeObject(result, f))
			if f.Required {
				rx := core.NewQRegExp()
				rx.SetPattern(".*")
				comboBox.SetValidator(gui.NewQRegExpValidator2(rx, comboBox))
			}
			field = comboBox
		case IntegerType:
			editor := widgets.NewQLineEdit(nil)
			editor.SetValidator(gui.NewQIntValidator2(-1000, 999999999, editor))
			editor.SetPlaceholderText(f.Label)
			editor.ConnectTextChanged(onChangeInt(result, f))
			if f.Required {
				rx := core.NewQRegExp()
				rx.SetPattern(".*")
				editor.SetValidator(gui.NewQRegExpValidator2(rx, editor))
			}
			field = editor
		case StringType:
			editor := widgets.NewQLineEdit(nil)
			editor.SetPlaceholderText(f.Label)
			editor.ConnectTextChanged(onChangeText(result, f))
			if f.Required {
				rx := core.NewQRegExp()
				rx.SetPattern(".*")
				editor.SetValidator(gui.NewQRegExpValidator2(rx, editor))
			}
			field = editor
		}
		if field != nil {
			layout.AddWidget(field, 0, 0)
		}
	}
	button := widgets.NewQPushButton2(buttonText, nil)
	button.ConnectClicked(func(c bool) {
		dialog.Accept()
	})
	layout.AddWidget(button, 0, core.Qt__AlignCenter)
	dialog.Exec()
}

func onChangeText(result map[string]interface{}, f Field) func(string string) {
	return func(text string) {
		if strings.TrimSpace(text) == "" {
			delete(result, f.Name)
		} else {
			result[f.Name] = text
		}
	}
}

func onChangeObject(result map[string]interface{}, f Field) func(string string) {
	return func(text string) {
		if strings.TrimSpace(text) == "" {
			delete(result, f.Name)
		} else {
			result[f.Name] = f.Options[text]
		}
	}
}

func onChangeInt(result map[string]interface{}, f Field) func(string string) {
	return func(text string) {
		if text == "" {
			delete(result, f.Name)
		} else {
			v, err := strconv.Atoi(text)
			if err == nil {
				result[f.Name] = v
			}
		}
	}
}
