package window

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (w Window) OpenInputWindow() string {
	println("Lets open window!")

	var inputEditor widget.Editor
	var submitButton widget.Clickable
	var input string

	w.gioWindow = new(app.Window)
	w.gioWindow.Option(app.Size(350, 250))
	w.gioWindow.Option(app.Title("🥴wacky_message"))

	theme := material.NewTheme()

	var ops op.Ops

	for {
		switch e := w.gioWindow.Event().(type) {
		case app.DestroyEvent:
			w.gioWindow = nil

			return input
		case app.FrameEvent:
			//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

			gtx := app.NewContext(&ops, e)

			title := material.H5(theme, "☠️Wacky Message🤪")
			title.Color = black
			title.Alignment = text.Middle

			paint.Fill(gtx.Ops, white)
			title.Layout(gtx)

			if submitButton.Clicked(gtx) {
				input = inputEditor.Text()
				inputEditor.SetText("")

				w.gioWindow.Perform(system.ActionClose)
			}

			layout.Center.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								return material.Editor(theme, &inputEditor, "Enter a message...").Layout(gtx)
							},
						),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								return material.Button(theme, &submitButton, "Modify and Insert").Layout(gtx)
							},
						),
					)
				},
			)

			e.Frame(gtx.Ops)
		}
	}
}
