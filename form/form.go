package form

import (
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/core"
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

func ShowForm(fields []Field, onClose func(map[string]interface{}), buttonText string, title string) {
	dialog := widgets.NewQDialog(nil, core.Qt__Dialog)
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
				for key, _ := range f.Options {
					comboBox.AddItem(key, nil)
				}
			}
			comboBox.ConnectCurrentIndexChanged2(func(text string) {
				result[f.Name] = f.Options[text]
			})
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
			editor.ConnectTextChanged(func(text string) {
				if text == "" {
					delete(result, f.Name)
				} else {
					v, err := strconv.Atoi(text)
					if err == nil {
						result[f.Name] = v
					}
				}
			})
			if f.Required {
				rx := core.NewQRegExp()
				rx.SetPattern(".*")
				editor.SetValidator(gui.NewQRegExpValidator2(rx, editor))
			}
			field = editor
		case StringType:
			editor := widgets.NewQLineEdit(nil)
			editor.SetPlaceholderText(f.Label)
			editor.ConnectTextChanged(onChange(result, f))
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
	dialog.Show()
}

func onChange(result map[string]interface{}, f Field) func(string string) {
	return func(text string) {
		if text == "" {
			delete(result, f.Name)
		} else {
			result[f.Name] = text
		}
	}
}
