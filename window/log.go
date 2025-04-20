package window

import (
	"image/color"
	"sync"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Log struct {
	gioWindow *app.Window
	logs      []string
	mu        sync.Mutex
}

func NewLog() *Log {
	return &Log{
		gioWindow: new(app.Window),
		logs:      []string{},
	}
}

func (lw *Log) Show() {
	lw.gioWindow = new(app.Window)
	lw.gioWindow.Option(app.Size(555, 250))
	lw.gioWindow.Option(app.Title("App log"))

	theme := material.NewTheme()

	var ops op.Ops
	list := widget.List{
		List: layout.List{Axis: layout.Vertical},
	}

	for {
		switch e := lw.gioWindow.Event().(type) {
		case app.DestroyEvent:
			return
		case app.FrameEvent:
			lw.mu.Lock()
			logs := lw.logs
			lw.mu.Unlock()

			gtx := app.NewContext(&ops, e)
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
			white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

			paint.Fill(gtx.Ops, black)

			list.Layout(
				gtx, len(logs), func(gtx layout.Context, i int) layout.Dimensions {
					label := material.Body1(theme, logs[i])
					label.Color = white
					return label.Layout(gtx)
				},
			)

			e.Frame(gtx.Ops)
		}
	}
}

func (lw *Log) AddLog(message string) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	lw.logs = append([]string{message}, lw.logs...)
}
