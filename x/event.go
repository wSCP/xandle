package x

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

type Eventer interface {
	Enqueue(xgb.Event, xgb.Error)
	Dequeue() Event
	ManageEventsOn(*xgb.Conn, chan struct{}, chan struct{}, chan struct{})
	Empty() bool
	Caller
	Ender
}

type Xevent int

const (
	UnknownEvent Xevent = iota
	MapRequest
	DestroyNotify
	UnmapNotify
	ClientMessage
	ConfigureRequest
	PropertyNotify
	EnterNotify
	MotionNotify
	FocusIn
	ScreenChange
	KeyPressEvent
	KeyReleaseEvent
	ButtonPressEvent
	ButtonReleaseEvent
)

type Event struct {
	Tag Xevent
	Evt xgb.Event
	Err xgb.Error
}

func newEvent(evt xgb.Event, err xgb.Error) Event {
	return Event{
		Tag: idEvent(evt),
		Evt: evt,
		Err: err,
	}
}

func idEvent(evt xgb.Event) Xevent {
	switch evt.(type) {
	case xproto.MapRequestEvent:
		return MapRequest
	case xproto.DestroyNotifyEvent:
		return DestroyNotify
	case xproto.UnmapNotifyEvent:
		return UnmapNotify
	case xproto.ClientMessageEvent:
		return ClientMessage
	case xproto.ConfigureRequestEvent:
		return ConfigureRequest
	case xproto.PropertyNotifyEvent:
		return PropertyNotify
	case xproto.EnterNotifyEvent:
		return EnterNotify
	case xproto.MotionNotifyEvent:
		return MotionNotify
	case xproto.FocusInEvent:
		return FocusIn
	case randr.ScreenChangeNotifyEvent:
		return ScreenChange
	case xproto.KeyPressEvent:
		return KeyPressEvent
	case xproto.KeyReleaseEvent:
		return KeyReleaseEvent
	case xproto.ButtonPressEvent:
		return ButtonPressEvent
	case xproto.ButtonReleaseEvent:
		return ButtonReleaseEvent
	}
	return UnknownEvent
}

type eventer struct {
	*log.Logger
	events []Event
	elk    *sync.RWMutex
	*caller
	*ender
}

func newEventer() *eventer {
	return &eventer{
		Logger: log.New(os.Stderr, "x:handle:eventer: ", log.Lmicroseconds|log.Llongfile),
		events: make([]Event, 0),
		elk:    &sync.RWMutex{},
		caller: newCaller(),
		ender:  newEnder(),
	}
}

func (e *eventer) Enqueue(evt xgb.Event, err xgb.Error) {
	e.elk.Lock()
	defer e.elk.Unlock()

	e.events = append(e.events, newEvent(evt, err))
}

func (e *eventer) Dequeue() Event {
	e.elk.Lock()
	defer e.elk.Unlock()

	ret := e.events[0]
	e.events = e.events[1:]
	return ret
}

func (e *eventer) ManageEventsOn(c *xgb.Conn, pre, post, quit chan struct{}) {
	for {
		if e.Ending() {
			if quit != nil {
				quit <- struct{}{}
			}
			break
		}

		e.read(c)
		e.process(pre, post)
	}
}

func (e *eventer) read(c *xgb.Conn) {
	ev, err := c.WaitForEvent()
	if ev == nil && err == nil {
		e.Fatal("eventer could not read an event or an error.")
	}
	e.Enqueue(ev, err)
}

func (e *eventer) process(pre, post chan struct{}) {
	for !e.Empty() {
		if e.Ending() {
			return
		}

		pre <- struct{}{}

		evt := e.Dequeue()

		v, err := evt.Evt, evt.Err

		if err != nil {
			e.Println(err.Error())
			post <- struct{}{}
			continue
		}

		if v == nil {
			e.Fatal("eventer.process expected an event but got nil.")
		}

		callErr := e.Call(evt)
		if callErr != nil {
			e.Println(callErr.Error())
		}

		post <- struct{}{}
	}
}

func (e *eventer) Empty() bool {
	e.elk.Lock()
	defer e.elk.Unlock()

	return len(e.events) == 0
}
