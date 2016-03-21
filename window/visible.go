package window

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Visible interface {
	Hide()
	Show()
}

type visible struct {
	c *xgb.Conn
	w xproto.Window
	r xproto.Window
}

var (
	RootEventMask uint32 = (xproto.EventMaskSubstructureNotify | xproto.EventMaskSubstructureRedirect)
	windowOff            = []uint32{RootEventMask, xproto.EventMaskSubstructureNotify} //uint32_t values_off[] = {ROOT_EVENT_MASK & ~XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY};
	windowOn             = []uint32{RootEventMask}
)

func setVisibility(v bool, c *xgb.Conn, w xproto.Window, root xproto.Window) {
	xproto.ChangeWindowAttributesChecked(c, root, xproto.CwEventMask, windowOff)
	if v {
		xproto.MapWindow(c, w)
	} else {
		xproto.UnmapWindow(c, w)
	}
	xproto.ChangeWindowAttributesChecked(c, root, xproto.CwEventMask, windowOn)
}

func (v *visible) Hide() {
	setVisibility(false, v.c, v.w, v.r)
}

func (v *visible) Show() {
	setVisibility(true, v.c, v.w, v.r)
}
