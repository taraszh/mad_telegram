package tray

type Tray interface {
	SetIcon(icon []byte)
	SetTooltip(tooltip string)
	AddMenu(title, tooltip string, onClick func())
	AddSeparator()
	Run(onReady func(), onExit func())
	Quit()
}
