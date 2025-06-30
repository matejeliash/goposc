package main

import (
	"fmt"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/matejeliash/goposc/internal/netinfo"
)

type Person struct {
	Name    string
	Age     int
	Address string
}

// func structToCanvasTexts(s any) fyne.CanvasObject {
// 	v := reflect.ValueOf(s)
// 	t := v.Type()

// 	content := container.NewVBox()

// 	for i := 0; i < v.NumField(); i++ {
// 		fieldName := t.Field(i).Name + ": "
// 		fieldValue := fmt.Sprintf("%v", v.Field(i).Interface())

// 		// Bold field name
// 		boldName := canvas.NewText(fieldName, nil)
// 		boldName.TextStyle = fyne.TextStyle{Bold: true}
// 		boldName.TextSize = 14

// 		// Normal value
// 		value := canvas.NewText(fieldValue, nil)
// 		value.TextSize = 14

// 		// Put both side by side
// 		row := container.NewHBox(boldName, value)
// 		content.Add(row)
// 	}

// 	return content
// }

// func structToCanvasTexts(s any) fyne.CanvasObject {
// 	v := reflect.ValueOf(s)
// 	t := v.Type()

// 	content := container.NewVBox()

// 	for i := 0; i < v.NumField(); i++ {
// 		fieldName := t.Field(i).Name + ": "
// 		fieldValue := v.Field(i)

// 		// Bold field name
// 		boldName := canvas.NewText(fieldName, nil)
// 		boldName.TextStyle = fyne.TextStyle{Bold: true}
// 		boldName.TextSize = 14

// 		if fieldValue.Kind() == reflect.Slice {
// 			// Create VBox for slice elements with indentation
// 			sliceContent := container.NewVBox()
// 			for j := 0; j < fieldValue.Len(); j++ {
// 				elemStr := fmt.Sprintf("%v", fieldValue.Index(j).Interface())
// 				// Indent slice elements by adding padding on the left using a spacer box or an empty text with spaces
// 				indentText := canvas.NewText("   "+elemStr, nil) // 3 spaces indent
// 				indentText.TextSize = 14
// 				sliceContent.Add(indentText)
// 			}

// 			// Put the field name on one row and slice elements in a VBox below
// 			content.Add(container.NewVBox(boldName, sliceContent))

// 		} else {
// 			// Normal value as text
// 			value := canvas.NewText(fmt.Sprintf("%v", fieldValue.Interface()), nil)
// 			value.TextSize = 14

// 			// Put both side by side
// 			row := container.NewHBox(boldName, value)
// 			content.Add(row)
// 		}
// 	}

// 	return content
// }

func structToCanvasTextsRecursive(v reflect.Value, indentLevel int) fyne.CanvasObject {
	// Dereference pointers if needed
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return canvas.NewText("<nil>", nil)
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		content := container.NewVBox()
		for i := 0; i < v.NumField(); i++ {
			fieldName := t.Field(i).Name + ": "

			// Bold field name with indentation
			indentSpaces := strings.Repeat(" ", indentLevel*3)
			boldName := canvas.NewText(indentSpaces+fieldName, nil)
			boldName.TextStyle = fyne.TextStyle{Bold: true}
			boldName.TextSize = 14

			fieldValue := v.Field(i)

			// Recursively get the value display with increased indent
			valueObj := structToCanvasTextsRecursive(fieldValue, indentLevel+1)

			// Combine field name + value vertically (name above value)
			content.Add(container.NewVBox(boldName, valueObj))
		}
		return content

	case reflect.Slice, reflect.Array:
		content := container.NewVBox()
		// Removed unused indentSpaces declaration here
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			// Recursively display element with increased indent
			elemText := structToCanvasTextsRecursive(elem, indentLevel+1)
			content.Add(elemText)
		}
		return content

	default:
		// Basic value: just display with indentation spaces
		indentSpaces := strings.Repeat(" ", indentLevel*3)
		text := canvas.NewText(indentSpaces+fmt.Sprintf("%v", v.Interface()), nil)
		text.TextSize = 14
		return text
	}
}

// Wrapper function for convenience
func structToCanvasTexts(s any) fyne.CanvasObject {
	return structToCanvasTextsRecursive(reflect.ValueOf(s), 0)
}

func main() {
	a := app.New()
	w := a.NewWindow("Struct Display")

	ni := &netinfo.NetworkInfo{}
	ni.GetInfos()
	fmt.Println(ni)

	w.SetContent(structToCanvasTexts(*ni))
	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}
