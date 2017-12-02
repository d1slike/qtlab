package form

import (
	"github.com/therecipe/qt/widgets"
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
	layout.SetSpacing(20)
	dialog.SetLayout(layout)
	dialog.SetWindowTitle(title)
	dialog.SetMinimumWidth(400)
	result := make(map[string]interface{})
	validations := make(map[string]bool)
	dialog.ConnectAccepted(func() {
		if len(result) > 0 {
			onClose(result)
		}
	})
	for _, f := range fields {
		var field widgets.QBoxLayout_ITF
		switch f.Type {
		case ObjectType:
			comboBox := widgets.NewQComboBox(dialog)
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
			if f.Required {
				validations[f.Name] = false
			}

			v := widgets.NewQVBoxLayout2(dialog)
			v.SetSpacing(5)
			label := widgets.NewQLabel2("Обязательно для заполнения", dialog, 0)
			label.SetStyleSheet("color: red; font-size: 7")
			label.SetVisible(f.Required)
			l := widgets.NewQLabel2(f.Label, dialog, 0)
			v.AddWidget(l, 0, 0)
			v.AddWidget(comboBox, 0, 0)
			v.AddWidget(label, 0, core.Qt__AlignRight)

			comboBox.ConnectCurrentIndexChanged2(onChangeObject(result, f, label, validations))
			field = v
		case IntegerType:
			editor := widgets.NewQLineEdit(dialog)
			editor.SetPlaceholderText(f.Label)
			if f.Required {
				validations[f.Name] = false
			}

			v := widgets.NewQVBoxLayout2(dialog)
			v.SetSpacing(5)
			label := widgets.NewQLabel2("Обязательно для заполнения", dialog, 0)
			label.SetStyleSheet("color: red; font-size: 7")
			label.SetVisible(f.Required)
			v.AddWidget(editor, 0, 0)
			v.AddWidget(label, 0, core.Qt__AlignRight)

			editor.ConnectTextChanged(onChangeInt(result, f, label, validations))
			field = v
		case StringType:
			editor := widgets.NewQLineEdit(dialog)
			editor.SetPlaceholderText(f.Label)
			if f.Required {
				validations[f.Name] = false
			}

			v := widgets.NewQVBoxLayout2(dialog)
			v.SetSpacing(5)
			label := widgets.NewQLabel2("Обязательно для заполнения", dialog, 0)
			label.SetStyleSheet("color: red; font-size: 7")
			label.SetVisible(f.Required)
			v.AddWidget(editor, 0, 0)
			v.AddWidget(label, 0, core.Qt__AlignRight)

			editor.ConnectTextChanged(onChangeText(result, f, label, validations))
			field = v
		}
		if field != nil {
			layout.AddLayout(field, 0)
		}
	}
	button := widgets.NewQPushButton2(buttonText, dialog)
	button.ConnectClicked(func(c bool) {
		allValid := true
		for _, valid := range validations {
			if !valid {
				allValid = false
				break
			}
		}
		if allValid {
			dialog.Accept()
		}
	})
	layout.AddWidget(button, 0, core.Qt__AlignCenter)
	dialog.Exec()
}

func onChangeText(result map[string]interface{}, f Field, label *widgets.QLabel, validations map[string]bool) func(string string) {
	return func(text string) {
		if strings.TrimSpace(text) == "" {
			delete(result, f.Name)
			label.SetVisible(f.Required)
			validations[f.Name] = !f.Required
		} else {
			label.SetVisible(false)
			result[f.Name] = text
			validations[f.Name] = true
		}
	}
}

func onChangeObject(result map[string]interface{}, f Field, label *widgets.QLabel, validations map[string]bool) func(string string) {
	return func(text string) {
		if strings.TrimSpace(text) == "" {
			delete(result, f.Name)
			label.SetVisible(f.Required)
			validations[f.Name] = !f.Required
		} else {
			label.SetVisible(false)
			result[f.Name] = f.Options[text]
			validations[f.Name] = true
		}
	}
}

func onChangeInt(result map[string]interface{}, f Field, label *widgets.QLabel, validations map[string]bool) func(string string) {
	return func(text string) {
		if text == "" {
			delete(result, f.Name)
			if f.Required {
				label.SetText("Обязательно для заполнения")
				label.SetVisible(true)
				validations[f.Name] = false
			}
		} else {
			v, err := strconv.Atoi(text)
			if err == nil {
				result[f.Name] = v
				label.SetVisible(false)
				validations[f.Name] = true
			} else {
				label.SetText("Введите целое число")
				label.SetVisible(true)
				validations[f.Name] = false
			}
		}
	}
}
