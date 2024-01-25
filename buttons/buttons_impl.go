package buttons

import (
	"time"

	gfx "github.com/ImVulcrum/Chess/gfxw"
)

type Butt struct {
	x            uint16
	y            uint16
	length       uint16
	height       uint16
	name         string
	r            uint8 //color of the button
	g            uint8
	b            uint8
	r_label      uint8 //color of the label
	g_label      uint8
	b_label      uint8
	label_offset uint16 //spacing between the left side of the button and the first char of the label (in pixels)
	font_size    int
	state        bool
	active       bool
}

func New(x uint16, y uint16, length uint16, height uint16, name string, re, gr, bl, r_label, g_label, b_label uint8, label_offset uint16, font_size int) *Butt {
	var b = new(Butt)
	b.x = x
	b.y = y
	b.length = length
	b.height = height
	b.name = name
	b.r = re
	b.g = gr
	b.b = bl
	b.r_label = r_label
	b.g_label = g_label
	b.b_label = b_label
	b.label_offset = label_offset
	b.font_size = font_size
	b.state = false
	b.active = true
	return b
}

func (b *Butt) Draw() {
	if b.active { //only draw if active
		gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
		gfx.Stiftfarbe(b.r, b.g, b.b)
		gfx.Vollrechteck((*b).x, (*b).y, (*b).length, (*b).height)
		gfx.Stiftfarbe(b.r_label, b.g_label, b.b_label)
		gfx.SchreibeFont((*b).x+(*b).label_offset, (*b).y+(*b).height/10, (*b).name)
	}
}

func (b *Butt) Is_Clicked(x, y uint16) bool { //returns true if a click on the button is executed and playes a animation if so
	if b.active && x >= b.x && x <= b.x+b.length && y >= b.y && y <= b.y+b.height {
		gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Transparenz(120)
		gfx.Vollrechteck((*b).x, (*b).y, (*b).length, (*b).height)
		time.Sleep(time.Duration(100) * time.Millisecond)
		gfx.Stiftfarbe(b.r, b.g, b.b)
		gfx.Transparenz(0)
		gfx.Vollrechteck((*b).x, (*b).y, (*b).length, (*b).height)
		gfx.Stiftfarbe(b.r_label, b.g_label, b.b_label)
		gfx.SchreibeFont((*b).x+(*b).label_offset, (*b).y+(*b).height/10, (*b).name)
		return true
	}
	return false
}

func (b *Butt) Give_State() bool {
	return b.state
}

func (b *Butt) Deactivate() {
	b.active = false
}

func (b *Butt) Activate() {
	b.active = true
}

func (b *Butt) Is_Active() bool {
	return b.active
}

func (b *Butt) Switch(re, gr, bl uint8) bool { //turns the button to a switch (with evervy execution of this function the button state is changed and it's color as well)
	if b.active && b.state {
		b.state = false
		gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
		gfx.Stiftfarbe(b.r, b.g, b.b)
		gfx.Vollrechteck(b.x, b.y, b.length, b.height)
		gfx.Stiftfarbe(b.r_label, b.g_label, b.b_label)
		gfx.SchreibeFont(b.x+b.label_offset, b.y+b.height/10, b.name)
	} else if b.active && !b.state {
		b.state = true
		gfx.SetzeFont("./resources/fonts/firamono.ttf", b.font_size)
		gfx.Stiftfarbe(re, gr, bl)
		gfx.Vollrechteck((*b).x, (*b).y, (*b).length, (*b).height)
		gfx.Stiftfarbe(b.r_label, b.g_label, b.b_label)
		gfx.SchreibeFont((*b).x+(*b).label_offset, (*b).y+(*b).height/10, (*b).name)
	}

	return b.state
}
