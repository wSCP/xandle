package monitor

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/wSCP/utils/branch"
	"github.com/thrisp/wSCP/utils/settings"
)

type Monitors struct {
	settings.Settings
	*branch.Branch
}

func NewMonitors(c *xgb.Conn, r xproto.Window, s *xproto.ScreenInfo) *Monitors {
	monitors := &Monitors{settings.New(), branch.New("monitors")}
	Initialize(monitors, c, r, s)
	return monitors
}

func Initialize(monitors *Monitors, c *xgb.Conn, r xproto.Window, s *xproto.ScreenInfo) {
	err := randr.Init(c)
	if err == nil {
		if err := Update(monitors, c, r); err == nil {
			randr.SelectInputChecked(c, r, randr.NotifyMaskScreenChange)
		} else {
			err = xinerama.Init(c)
			if err == nil {
				xia, err := xinerama.IsActive(c).Reply()
				if xia != nil && err == nil {
					xsq, _ := xinerama.QueryScreens(c).Reply()
					xsi := xsq.ScreenInfo
					for i := 0; i < len(xsi); i++ {
						info := xsi[i]
						rect := xproto.Rectangle{info.XOrg, info.YOrg, info.Width, info.Height}
						nm := New(uint32(i), fmt.Sprintf("XMonitor%d", i), c, r, rect, monitors.Settings)
						monitors.PushBack(nm)
					}
				} else {
					rect := xproto.Rectangle{0, 0, s.WidthInPixels, s.HeightInPixels}
					nm := New(1, "SCREEN", c, r, rect, monitors.Settings)
					monitors.PushBack(nm)
				}
			}
		}
	}
}

//func Focus(monitors *branch.Branch, mon Monitor) {
//	focused := Focused(monitors)
//	focused.Set("focused", false)
//	mon.Focus()
//}

//func FocusSelect(sel []selector.Selector, monitors *branch.Branch) {
//	selected := Select(sel, monitors)
//	if selected != nil {
//		Focus(monitors, selected)
//	}
//}
