package systray_adapter

import (
	"github.com/getlantern/systray"
)

type SystrayAdapter struct{}

func (s *SystrayAdapter) SetIcon(icon []byte) {
	systray.SetIcon(icon)
}

func (s *SystrayAdapter) SetTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}

func (s *SystrayAdapter) AddMenu(title string, tooltip string, onClick func()) {
	item := systray.AddMenuItem(title, tooltip)
	go func() {
		for {
			select {
			case <-item.ClickedCh:
				println("Menu item clicked:", title)
				onClick()
			}
		}
	}()
}

func (s *SystrayAdapter) AddSeparator() {
	systray.AddSeparator()
}

func (s *SystrayAdapter) Run(onReady func(), onExit func()) {
	systray.Run(onReady, onExit)
}

func (s *SystrayAdapter) Quit() {
	systray.Quit()
}
