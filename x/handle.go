package x

import (
	"log"
	"os"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/wSCP/xandle/x/atomic"
	"github.com/wSCP/xandle/x/icccm"
)

type Manager interface {
	Manage(*xgb.Conn, chan struct{}, chan struct{}, chan struct{})
}

type Handle interface {
	Connection
	Inform
	Eventer
	Manager
	atomic.Atomic
	icccm.ICCCM
}

type handle struct {
	c *configuration
	*log.Logger
	*connection
	*inform
	*eventer
	atomic.Atomic
	icccm.ICCCM
}

func New(display string, cnf ...Config) (Handle, error) {
	c, err := xgb.NewConnDisplay(display)
	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)
	root := screen.Root

	h := &handle{
		c:          newConfiguration(cnf...),
		Logger:     log.New(os.Stderr, "x:handle: ", log.Lmicroseconds|log.Llongfile),
		connection: newConnection(c),
		inform:     newInform(setup, screen, root),
		eventer:    newEventer(),
		Atomic:     atomic.New(c),
	}

	i := icccm.New(h.Conn(), h.Root(), h.Atomic)
	h.ICCCM = i

	return h, err
}

func (h *handle) Manage(c *xgb.Conn, pre, post, quit chan struct{}) {
	if !h.c.Configured() {
		h.c.Configure(h)
	}
	h.ManageEventsOn(c, pre, post, quit)
}
