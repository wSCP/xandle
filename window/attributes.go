package window

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Attributes interface {
	MapState() int
	Mapped() bool
	Viewable() bool
}

type attributes struct {
	c *xgb.Conn
	w xproto.Window
}

func NewAttributes(c *xgb.Conn, w xproto.Window) Attributes {
	return &attributes{c, w}
}

func (a *attributes) getAttributes() *xproto.GetWindowAttributesReply {
	if ret, err := xproto.GetWindowAttributes(a.c, a.w).Reply(); err == nil {
		return ret
	}
	return nil
}

func (a *attributes) MapState() int {
	ms := a.getAttributes()
	if ms != nil {
		switch ms.MapState {
		case xproto.MapStateUnmapped:
			return 0
		case xproto.MapStateUnviewable:
			return 1
		case xproto.MapStateViewable:
			return 2
		}
	}
	return -1
}

func (a *attributes) Mapped() bool {
	m := a.MapState()
	if m == 1 || m == 2 {
		return true
	}
	return false
}

func (a *attributes) Viewable() bool {
	m := a.MapState()
	if m == 2 {
		return true
	}
	return false
}
