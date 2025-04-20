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

func (w Window) OpenHotKeySettings(windowClassMap map[string]string, selectedWindows []string) []string {
	println("Lets open hotkey setting window!")

	var checkboxes []widget.Bool
	var keys []string
	var values []string

	for key, value := range windowClassMap {
		checkbox := widget.Bool{}
		for _, selected := range selectedWindows {
			if selected == value {
				checkbox.Value = true
				break
			}
		}

		checkboxes = append(checkboxes, checkbox)
		keys = append(keys, key)
		values = append(values, value)
	}

	var submitButton widget.Clickable
	var selectedKeys []string

	w.gioWindow = new(app.Window)
	w.gioWindow.Option(app.Size(400, 600))
	w.gioWindow.Option(app.Title("Hotkey settings"))

	theme := material.NewTheme()

	var ops op.Ops

	for {
		switch e := w.gioWindow.Event().(type) {
		case app.DestroyEvent:
			w.gioWindow = nil
			return selectedKeys
		case app.FrameEvent:
			white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

			gtx := app.NewContext(&ops, e)

			paint.Fill(gtx.Ops, black)

			if submitButton.Clicked(gtx) {
				println("Submit button clicked")
				selectedKeys = nil

				for i, checkbox := range checkboxes {
					if checkbox.Value {
						println("Selected:", values[i])
						selectedKeys = append(selectedKeys, values[i])
					}
				}

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
								title := material.H5(theme, "Choose windows where hotkey will be active")
								title.Color = white
								title.Alignment = text.Middle
								return title.Layout(gtx)
							},
						),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(
									gtx,
									func() []layout.FlexChild {
										var children []layout.FlexChild
										for i, _ := range values {
											index := i
											children = append(
												children, layout.Rigid(
													func(gtx layout.Context) layout.Dimensions {
														return layout.Flex{
															Axis: layout.Horizontal,
														}.Layout(
															gtx,
															layout.Rigid(
																material.CheckBox(
																	theme,
																	&checkboxes[index],
																	"",
																).Layout,
															),
															layout.Rigid(
																func(gtx layout.Context) layout.Dimensions {
																	label := material.Body1(theme, keys[index])
																	label.Color = white
																	return label.Layout(gtx)
																},
															),
														)
													},
												),
											)
										}
										return children
									}()...,
								)
							},
						),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								return material.Button(theme, &submitButton, "Submit").Layout(gtx)
							},
						),
					)
				},
			)

			e.Frame(gtx.Ops)
		}
	}
}
