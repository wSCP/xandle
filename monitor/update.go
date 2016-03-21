package monitor

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

var UpdateError = Xrror("Monitors update error: %s").Out

func Update(monitors *Monitors, c *xgb.Conn, r xproto.Window) error {
	sres, err := randr.GetScreenResources(c, r).Reply()
	if err != nil {
		return err
	}

	var rr []*randr.GetOutputInfoReply
	for _, o := range sres.Outputs {
		rp, _ := randr.GetOutputInfo(c, o, xproto.TimeCurrentTime).Reply()
		rr = append(rr, rp)
	}

	curr := monitors.Front()
	for curr != nil {
		m := curr.Value.(Monitor)
		m.Set("wired", false)
		curr = curr.Next()
	}

	for i, info := range rr {
		if info != nil {
			if info.Crtc != 0 {
				ir, _ := randr.GetCrtcInfo(c, info.Crtc, xproto.TimeCurrentTime).Reply()
				if ir != nil {
					rect := xproto.Rectangle{ir.X, ir.Y, ir.Width, ir.Height}
					m := fromId(monitors, uint32(sres.Outputs[i]))
					if m != nil {
						m.SetRectangle(rect)
						m.UpdateRoot()
						m.Set("wired", true)
					} else {
						nm := New(uint32(sres.Outputs[i]), string(info.Name), c, r, rect, monitors.Settings)
						monitors.PushBack(nm)
					}
				}
			} else if !monitors.Bool("RemoveDisabled") && info.Connection != randr.ConnectionDisconnected {
				m := fromId(monitors, uint32(sres.Outputs[i]))
				if m != nil {
					m.Set("wired", true)
				}
			}
		}
	}

	gpo, _ := randr.GetOutputPrimary(c, r).Reply()
	if gpo != nil {
		pm := fromId(monitors, uint32(gpo.Output))
		if pm != nil {
			pm.Set("primary", true)
			pm.Set("focused", true)
		}
	}

	if monitors.Bool("MergeOverlapping") {
		mergeOverlapping(monitors)
	}

	if monitors.Bool("RemoveUnplugged") {
		removeUnplugged(monitors)
	}

	ml := monitors.Len()
	if ml > 0 {
		return nil
	} else {
		return UpdateError("Monitors update returned length %d monitors", ml)
	}
}

type UpdateMonitor func(Monitor) error

func update(monitors *Monitors, fn UpdateMonitor) error {
	curr := monitors.Front()
	for curr != nil {
		mon := curr.Value.(Monitor)
		if err := fn(mon); err != nil {
			return err
		}
		curr = curr.Next()
	}
	return nil
}

func contains(a, b xproto.Rectangle) bool {
	return (a.X <= b.X &&
		(a.X+int16(a.Width)) >= (b.X+int16(b.Width)) &&
		a.Y <= b.Y && (a.Y+int16(a.Height)) >= (b.Y+int16(b.Height)))
}

func mergeOverlapping(monitors *Monitors) {
	var mon1, mon2 Monitor
	m := monitors.Front()
	for m != nil {
		mon1 = m.Value.(Monitor)
		if mon1.Wired() {
			mm := monitors.Front()
			for mm != nil {
				mon2 = mm.Value.(Monitor)
				if mon1.Id() != mon2.Id() && mon2.Wired() && contains(mon1.Rectangle(), mon2.Rectangle()) {
					mon1.Merge(mon2)
				}
				mm = mm.Next()
			}
		}
		m = m.Next()
	}
}

func removeUnplugged(monitors *Monitors) {
	fn := func(mon Monitor) error {
		if !mon.Wired() {
			focused := Focused(monitors)
			focused.Merge(mon)
		}
		return nil
	}
	update(monitors, fn)
}
