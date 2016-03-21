package window

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Stack interface {
	Raise()
	Lower()
	Above(xproto.Window)
	Below(xproto.Window)
}

type stack struct {
	c *xgb.Conn
	w xproto.Window
}

func (s *stack) Raise() {
	xproto.ConfigureWindowChecked(s.c, s.w, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeAbove})
}

func (s *stack) Lower() {
	xproto.ConfigureWindowChecked(s.c, s.w, xproto.ConfigWindowStackMode, []uint32{xproto.StackModeBelow})
}

func (s *stack) stack(o xproto.Window, mode uint32) {
	xproto.ConfigureWindowChecked(
		s.c,
		s.w,
		(xproto.ConfigWindowSibling | xproto.ConfigWindowStackMode),
		[]uint32{uint32(o), mode},
	)
}

func (s *stack) Above(o xproto.Window) {
	s.stack(o, xproto.StackModeAbove)
}

func (s *stack) Below(o xproto.Window) {
	s.stack(o, xproto.StackModeBelow)
}
