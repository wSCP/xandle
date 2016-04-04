package x

import "github.com/BurntSushi/xgb/xproto"

type Inform interface {
	Setup() *xproto.SetupInfo
	Screen() *xproto.ScreenInfo
	Root() xproto.Window
}

type inform struct {
	setup  *xproto.SetupInfo
	screen *xproto.ScreenInfo
	root   xproto.Window
}

func newInform(si *xproto.SetupInfo, scr *xproto.ScreenInfo, r xproto.Window) *inform {
	return &inform{si, scr, r}
}

func (i *inform) Setup() *xproto.SetupInfo {
	return i.setup
}

func (i *inform) Screen() *xproto.ScreenInfo {
	return i.screen
}

func (i *inform) Root() xproto.Window {
	return i.root
}
