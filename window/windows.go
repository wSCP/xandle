package window

import (
	"strconv"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Windows interface {
	ListWindows() ([]xproto.Window, error)
	WindowExists(string) (xproto.Window, bool)
	RootWindow() Window
}

type windows struct {
	c  *xgb.Conn
	r  xproto.Window
	rt Window
}

func NewWindows(c *xgb.Conn, r xproto.Window) Windows {
	return &windows{c, r, New(c, r, r)}
}

func stringWindow(w string) xproto.Window {
	wid, err := strconv.ParseUint(w, 10, 32)
	if err != nil {
		return 0
	}
	return xproto.Window(wid)
}

func (w *windows) ListWindows() ([]xproto.Window, error) {
	q, err := xproto.QueryTreeUnchecked(w.c, w.r).Reply()
	if err != nil {
		return nil, err
	}
	return q.Children, nil
}

func (w *windows) WindowExists(requested string) (xproto.Window, bool) {
	wl, err := w.ListWindows()
	if err != nil {
		return 0, false
	}
	req := stringWindow(requested)
	for _, wid := range wl {
		if wid == req {
			return wid, true
		}
	}
	return 0, false
}

func (w *windows) RootWindow() Window {
	return w.rt
}
