package chocolate

import (
	"math"
)

type OverlayPosition int

const (
	CENTER OverlayPosition = iota
	START
	END
)

type Overlay struct {
	Chocolate

	enabled bool
	zindex  int

	width  float64
	height float64

	xpos OverlayPosition
	ypos OverlayPosition
}

func (o *Overlay) Enable()         { o.enabled = true }
func (o *Overlay) Disable()        { o.enabled = false }
func (o *Overlay) SetZIndex(v int) { o.zindex = v }
func (o *Overlay) SetPosition(pos ...OverlayPosition) {
	if len(pos) >= 1 {
		o.xpos = pos[0]
	}
	if len(pos) >= 2 {
		o.ypos = pos[1]
	}
}

func (o *Overlay) Resize(width, height int) {
	w := width
	h := height
	ow, owm := parseSize(o.width, 3)
	oh, ohm := parseSize(o.height, 3)
	// fmt.Printf("%d %d\n", ow, owm)
	// fmt.Printf("%d %d\n", oh, ohm)
	if ow > 0 && (ow+owm) < w {
		w = ow
	}
	if oh > 0 && (oh+ohm) < h {
		h = oh
	}
	if ow < 0 && (((ow * -1) + owm) < w) {
		w = w + ow
	}
	if oh < 0 && (((oh * -1) + ohm) < h) {
		h = h + oh
	}
	if ow == 0 {
		w = w - owm
	}
	if oh == 0 {
		h = h - ohm
	}
	o.Chocolate.Resize(w, h)
}

func newOverlay(choc *Chocolate, zindex int, width float64, height float64, pos ...OverlayPosition) *Overlay {
	ret := &Overlay{
		Chocolate: *choc,
		enabled:   false,
		zindex:    zindex,
		width:     width,
		height:    height,
	}
	ret.SetPosition(pos...)

	return ret
}

func (o *Overlay) calcPosition(pw, ph, w, h int) (int, int) {
	x := 0
	y := 0
	_, owm := parseSize(o.width, 3)
	_, ohm := parseSize(o.height, 3)
	switch o.xpos {
	case CENTER:
		x = (pw-w)/2 + owm
	case START:
		x = 0 + owm
	case END:
		x = pw - w - owm
	}
	switch o.ypos {
	case CENTER:
		y = (ph-h)/2 + ohm
	case START:
		y = 0 + ohm
	case END:
		y = ph - h - ohm
	}

	return x, y
}

func parseSize(v float64, precision uint) (int, int) {
	ratio := math.Pow(10, float64(precision))
	rf := math.Round(v*ratio) / ratio
	first := int(rf)
	second := int((v - float64(first)) * ratio)
	if second < 0 {
		second = second * -1
	}

	return first, second
}
