package xandle

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/thrisp/wSCP/xandle/atomic"
	"github.com/thrisp/wSCP/xandle/icccm"
	"github.com/thrisp/wSCP/xandle/monitor"
	"github.com/thrisp/wSCP/xandle/window"
)

type Xandle interface {
	Connectr
	Informr
	Eventr
	Callr
	atomic.Atomic
	icccm.ICCCM
	window.Windows
	Monitors() *monitor.Monitors
}

type Connectr interface {
	Conn() *xgb.Conn
}

type Informr interface {
	Setup() *xproto.SetupInfo
	Screen() *xproto.ScreenInfo
	Root() xproto.Window
}

type Eventr interface {
	Enqueue(xgb.Event, xgb.Error)
	Dequeue() (xgb.Event, xgb.Error)
	Handle(chan struct{}, chan struct{}, chan struct{})
	Empty() bool
	Endr
}

type Endr interface {
	End()
	Ending() bool
}

type Call func(xgb.Event) error

type Callr interface {
	SetEventFn(string, Call)
	Call(string, xgb.Event)
}

type xandle struct {
	*log.Logger
	conn    *xgb.Conn
	setup   *xproto.SetupInfo
	screen  *xproto.ScreenInfo
	root    xproto.Window
	events  []evnt
	evtsLck *sync.RWMutex
	call    map[string]Call
	callLck *sync.RWMutex
	end     bool
	atomic.Atomic
	icccm.ICCCM
	window.Windows
	monitors *monitor.Monitors
}

//func New(display string, ewhm []string, logr *log.Logger) (Xandle, error) {
func New(display string) (Xandle, error) {
	c, err := xgb.NewConnDisplay(display)
	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)

	x := &xandle{
		Logger:  log.New(os.Stderr, "xandle: ", log.Lshortfile|log.Lmicroseconds),
		conn:    c,
		setup:   setup,
		screen:  screen,
		root:    screen.Root,
		events:  make([]evnt, 0, 1000),
		evtsLck: &sync.RWMutex{},
		call:    make(map[string]Call),
		callLck: &sync.RWMutex{},
	}

	x.Atomic = atomic.New(x.conn)
	//x.Atomic.Atom("WM_DELETE_WINDOW")
	//x.Atomic.Atom("WM_TAKE_FOCUS")

	i := icccm.New(x.conn, x.root, x.Atomic)
	//i.SupportedSet(ewhm)
	//if err != nil {
	//	return nil, err
	//}
	//x.Ewmh.Set("string name", x.root, x.meta)
	x.ICCCM = i

	x.Windows = window.NewWindows(x.conn, x.root)

	x.monitors = monitor.NewMonitors(x.conn, x.root, x.screen)
	return x, nil
}

func (x *xandle) Conn() *xgb.Conn {
	return x.conn
}

func (x *xandle) Setup() *xproto.SetupInfo {
	return x.setup
}

func (x *xandle) Screen() *xproto.ScreenInfo {
	return x.screen
}

func (x *xandle) Root() xproto.Window {
	return x.root
}

func (x *xandle) Empty() bool {
	x.evtsLck.Lock()
	defer x.evtsLck.Unlock()

	return len(x.events) == 0
}

func (x *xandle) End() {
	x.end = true
}

func (x *xandle) Ending() bool {
	return x.end
}

type evnt struct {
	evt xgb.Event
	err xgb.Error
}

func (x *xandle) Enqueue(evt xgb.Event, err xgb.Error) {
	x.evtsLck.Lock()
	defer x.evtsLck.Unlock()

	x.events = append(x.events, evnt{
		evt: evt,
		err: err,
	})
}

func (x *xandle) Dequeue() (xgb.Event, xgb.Error) {
	x.evtsLck.Lock()
	defer x.evtsLck.Unlock()

	e := x.events[0]
	x.events = x.events[1:]
	return e.evt, e.err
}

func (x *xandle) Handle(pre, post, quit chan struct{}) {
	for {
		if x.Ending() {
			if quit != nil {
				quit <- struct{}{}
			}
			break
		}

		x.read()

		x.process(pre, post)
	}
}

func (x *xandle) read() {
	ev, err := x.Conn().WaitForEvent()
	if ev == nil && err == nil {
		x.Fatal("BUG: Could not read an event or an error.")
	}
	x.Enqueue(ev, err)
}

func (x *xandle) process(pre, post chan struct{}) {
	for !x.Empty() {
		if x.Ending() {
			return
		}

		pre <- struct{}{}

		ev, err := x.Dequeue()

		if err != nil {
			x.Println(err.Error())
			post <- struct{}{}
			continue
		}

		if ev == nil {
			x.Fatal("BUG: Expected an event but got nil.")
		}

		var tag string
		switch ev.(type) {
		case xproto.MapRequestEvent:
			tag = "MapRequest"
		case xproto.DestroyNotifyEvent:
			tag = "DestroyNotify"
		case xproto.UnmapNotifyEvent:
			tag = "UnmapNotify"
		case xproto.ClientMessageEvent:
			tag = "ClientMessage"
		case xproto.ConfigureRequestEvent:
			tag = "ConfigureRequest"
		case xproto.PropertyNotifyEvent:
			tag = "PropertyNotify"
		case xproto.EnterNotifyEvent:
			tag = "EnterNotify"
		case xproto.MotionNotifyEvent:
			tag = "MotionNotify"
		case xproto.FocusInEvent:
			tag = "FocusIn"
		case randr.ScreenChangeNotifyEvent:
			tag = "ScreenChange"
		}

		if tag != "" {
			x.Call(tag, ev)
		}

		post <- struct{}{}
	}
}

func (x *xandle) SetEventFn(tag string, fn Call) {
	x.callLck.Lock()
	defer x.callLck.Unlock()
	x.call[tag] = fn
}

func (x *xandle) Call(tag string, evt xgb.Event) {
	x.callLck.Lock()
	defer x.callLck.Unlock()
	if fn, ok := x.call[tag]; ok {
		err := fn(evt)
		if err != nil {
			x.Println("ERROR: %s", err.Error())
		}
	}
}

func (x *xandle) Monitors() *monitor.Monitors {
	return x.monitors
}
