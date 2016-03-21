package window

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Location interface {
	GetXY() (uint32, uint32)
	SetXY(uint32, uint32)
	GetX() uint32
	SetX(uint32)
	GetY() uint32
	SetY(uint32)
}

type Dimension interface {
	GetHeightWidth() (uint32, uint32)
	SetHeightWidth(uint32, uint32)
	GetWidth() uint32
	SetWidth(uint32)
	GetHeight() uint32
	SetHeight(uint32)
}

type Border interface {
	GetBorderWidth() uint32
	SetBorderWidth(uint32)
}

type Teleporter interface {
	MoveResize(uint32, uint32, uint32, uint32)
}

type Geometry interface {
	Location
	Dimension
	Teleporter
	Border
}

type geometry struct {
	c *xgb.Conn
	w xproto.Window
}

func NewGeometry(c *xgb.Conn, w xproto.Window) Geometry {
	return &geometry{c, w}
}

func (g *geometry) getGeometry() *xproto.GetGeometryReply {
	resp, err := xproto.GetGeometry(g.c, xproto.Drawable(g.w)).Reply()
	if err == nil {
		return resp
	}
	return nil
}

func (g *geometry) GetXY() (uint32, uint32) {
	geo := g.getGeometry()
	return uint32(geo.X), uint32(geo.Y)
}

func (g *geometry) SetXY(x, y uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowX|xproto.ConfigWindowY, []uint32{x, y})
}

func (g *geometry) GetX() uint32 {
	return uint32(g.getGeometry().X)
}

func (g *geometry) SetX(x uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowX, []uint32{x})
}

func (g *geometry) GetY() uint32 {
	return uint32(g.getGeometry().Y)
}

func (g *geometry) SetY(y uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowY, []uint32{y})
}

func (g *geometry) GetHeightWidth() (uint32, uint32) {
	geo := g.getGeometry()
	return uint32(geo.Height), uint32(geo.Width)
}

func (g *geometry) SetHeightWidth(hght, wdth uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowHeight|xproto.ConfigWindowWidth, []uint32{hght, wdth})
}

func (g *geometry) GetWidth() uint32 {
	return uint32(g.getGeometry().Width)
}

func (g *geometry) SetWidth(wdth uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowHeight, []uint32{wdth})
}

func (g *geometry) GetHeight() uint32 {
	return uint32(g.getGeometry().Height)
}

func (g *geometry) SetHeight(hght uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowWidth, []uint32{hght})
}

func (g *geometry) MoveResize(x, y, hght, wdth uint32) {
	g.SetXY(x, y)
	g.SetHeightWidth(hght, wdth)
}

func (g *geometry) GetBorderWidth() uint32 {
	return uint32(g.getGeometry().BorderWidth)
}

func (g *geometry) SetBorderWidth(bw uint32) {
	xproto.ConfigureWindowChecked(g.c, g.w, xproto.ConfigWindowBorderWidth, []uint32{bw})
}
