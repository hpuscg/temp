package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"github.com/flopp/go-findfont"
)

func main() {
	Init()
	a := app.New()
	w := a.NewWindow("进场自动化程序V1.0")

	MainShow(w)
	w.Resize(fyne.Size{Width: 400, Height: 80})
	w.CenterOnScreen()

	w.ShowAndRun()

	if err := os.Unsetenv("FYNE_FONT"); err != nil {
		return
	}
}

func Init() {
	fontPaths := findfont.List()
	// fmt.Println("=========")
	for _, fontPath := range fontPaths {
		// fmt.Println("============", fontPath)
		if strings.Contains(fontPath, "Songti.ttc") {
			fmt.Println("=========")
			if err := os.Setenv("FYNE_FONT", fontPath); err != nil {
				return
			}
			break
		}
	}
}

func MainShow(w fyne.Window) {
	title := widget.NewLabel("进场自动化程序")
	hello := widget.NewLabel("文件夹路径：")
	entry1 := widget.NewEntry()
	dia1 := widget.NewButton("打开", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				fmt.Println("Cancelled")
				return
			}
			entry1.SetText(reader.URI().String())
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})
	label2 := widget.NewLabel("切面方式：")
	text := widget.NewMultiLineEntry()
	text.Disable()
	labelLast := widget.NewLabel("                   ")
	combox1 := widget.NewSelect([]string{"最大切面值", "固定倾角切面"}, func(s string) { fmt.Println("selected", s) })
	label3 := widget.NewLabel("极化方式：")
	combox1.Resize(fyne.Size{Width: 200, Height: 90})
	combox2 := widget.NewSelect([]string{"45极化", "H/V极化"}, func(s string) { fmt.Println("selected", s) })
	label4 := widget.NewLabel("结果文件夹：")
	entry2 := widget.NewEntry()
	dia2 := widget.NewButton("打开", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if list == nil {
				fmt.Println("Cancelled")
				return
			}
			entry2.SetText(list.String())
		}, w)
	})
	combox1.SetSelectedIndex(0)
	combox2.SetSelectedIndex(0)
	bt3 := widget.NewButton("生成脚本", func() {
		if entry1.Text != "" && entry2.Text != "" {
			text.SetText("")
			text.Refresh()
			generateTxt(entry1.Text, entry2.Text, combox2.Selected, w)
			text.SetText("ok")
			text.Refresh()
		} else {
			dialog.ShowError(errors.New("failed"), w)
		}
	})
	bt4 := widget.NewButton("汇总结果", func() {
		fmt.Println(entry2.Text)
		if entry2.Text != "" {
			bt2(entry2.Text, combox1.Selected, combox2.Selected, text)
		} else {
			dialog.ShowError(errors.New("failed"), w)
		}
	})
	head := container.NewCenter(title)
	v1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), hello, dia1, entry1)
	v2 := container.NewHBox(label2, combox1)
	v3 := container.NewHBox(label3, combox2)
	v4 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), label4, dia2, entry2)
	v5 := container.NewHBox(bt3, bt4)
	v5Center := container.NewCenter(v5)
	ctnt := container.NewVBox(head, v1, v2, v3, v4, v5Center, text, labelLast)
	w.SetContent(ctnt)
}

func generateTxt(string, string, string, fyne.Window) string {
	return ""
}

func bt2(string, string, string, *widget.Entry) {
	fmt.Println("no")
}
