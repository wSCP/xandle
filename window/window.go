package window

import (
	"strconv"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Window interface {
	Id() xproto.Window
	IdString() string
	Attributes
	Geometry
	Visible
	Stack
}

type window struct {
	w xproto.Window
	Attributes
	Geometry
	Visible
	Stack
}

func New(c *xgb.Conn, w xproto.Window, r xproto.Window) Window {
	return &window{
		w,
		NewAttributes(c, w),
		NewGeometry(c, w),
		&visible{c, w, r},
		&stack{c, w},
	}
}

func (w *window) Id() xproto.Window {
	return w.w
}

func (w *window) IdString() string {
	return strconv.FormatUint(uint64(w.w), 10)
}

/*
//Conn() *xgb.Conn
//XRoot() xproto.Window
//XWindow() xproto.Window
//Close()
//Kill()

func (w *window) Conn() *xgb.Conn {
	return w.c
}

func (w *window) XRoot() xproto.Window {
	return w.r
}

func (w *window) Close() {
	//send_client_message(w.Window, ewmh->WM_PROTOCOLS, WM_DELETE_WINDOW);
}

func (w *window) Kill() {
	xproto.KillClientChecked(w.c, uint32(w.w))
}
*/
