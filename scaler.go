package chocolate

type contentAddHandlerFct func(int)

type fixedScaler struct {
	value  int
	addFct contentAddHandlerFct
}

func (s *fixedScaler) finalSize(size int) int { return s.value }
func (s *fixedScaler) setCotentSize(size int) int {
	if s.addFct != nil {
		s.addFct(s.value)
	}
	return s.value
}

func (s *fixedScaler) getSize() int      { return s.value }
func (s *fixedScaler) setSize(value int) { s.value = value }

type dynamicScaler struct {
	addFct contentAddHandlerFct
}

func (s *dynamicScaler) finalSize(size int) int { return size }
func (s *dynamicScaler) setCotentSize(size int) int {
	if s.addFct != nil {
		s.addFct(size)
	}
	return size
}

type (
	partsTakeHandlerFct func(int) int
	partsAddHandlerFct  func(int)
)

func (s *dynamicScaler) getSize() int      { return 0 }
func (s *dynamicScaler) setSize(value int) {}

type parentScaler struct {
	value   int
	addFct  partsAddHandlerFct
	takeFct partsTakeHandlerFct
}

func (s *parentScaler) finalSize(size int) int {
	if s.takeFct != nil {
		return s.takeFct(s.value)
	}
	return 0
}

func (s *parentScaler) setCotentSize(size int) int {
	if s.addFct != nil {
		s.addFct(s.value)
	}
	return 0
}

func (s *parentScaler) getSize() int      { return s.value }
func (s *parentScaler) setSize(value int) { s.value = value }

type fixedCreator struct {
	value int
}

func (c *fixedCreator) create(adder ContentAdder) Scaler {
	if adder == nil {
		return &fixedScaler{
			value: c.value,
		}
	}
	return &fixedScaler{
		value:  c.value,
		addFct: adder.AddContent,
	}
}

type dynamicCreator struct{}

func (c *dynamicCreator) create(adder ContentAdder) Scaler {
	if adder == nil {
		return &dynamicScaler{}
	}
	return &dynamicScaler{
		addFct: adder.AddContent,
	}
}

type parentCreator struct {
	value int
}

func (c *parentCreator) create(adder ContentAdder) Scaler {
	if adder == nil {
		return &parentScaler{
			value: c.value,
		}
	}
	return &parentScaler{
		value:   c.value,
		addFct:  adder.AddParts,
		takeFct: adder.TakeParts,
	}
}

type ContentAdder interface {
	AddContent(size int)
	AddParts(size int)
	TakeParts(size int) int
	GetMaxSize() int
}

type Scaler interface {
	finalSize(size int) int
	setCotentSize(size int) int
	getSize() int
	setSize(value int)
}

type createScaler interface {
	create(ContentAdder) Scaler
}

type Sizer interface {
	ContentAdder
	GetSize() int
	GetMaxSize() int
	GetContentSize() int
	SetSize(int)
	SetMaxSize(int)
	SetContentSize(int)
}

type scaler struct {
	xp ContentAdder
	yp ContentAdder
	x  Scaler
	y  Scaler

	xc createScaler
	yc createScaler

	xs Sizer
	ys Sizer
}

func (s *scaler) getWidth() int  { return s.xs.GetSize() }
func (s *scaler) getHeight() int { return s.ys.GetSize() }

func (s *scaler) setWidth(v int)  { s.xs.SetSize(v) }
func (s *scaler) setHeight(v int) { s.ys.SetSize(v) }

func (s *scaler) setMaxWidth(v int)  { s.xs.SetMaxSize(v) }
func (s *scaler) setMaxHeight(v int) { s.ys.SetMaxSize(v) }

func (s *scaler) getParentMaxWidth() int {
	if s.xp != nil {
		return s.xp.GetMaxSize()
	}
	return 0
}

func (s *scaler) getParentMaxHeight() int {
	if s.yp != nil {
		return s.yp.GetMaxSize()
	}
	return 0
}

func (s *scaler) setContentSize(width, height int) {
	if s.x != nil {
		s.x.setCotentSize(width)
	}
	if s.y != nil {
		s.y.setCotentSize(height)
	}
	s.xs.SetContentSize(width)
	s.ys.SetContentSize(height)
}

func (s *scaler) finalizeSize(xmargin, ymargin int) (int, int) {
	w := s.getWidth()
	h := s.getHeight()

	w = s.x.finalSize(w)
	h = s.y.finalSize(h)

	w -= xmargin
	h -= ymargin

	s.setWidth(w)
	s.setHeight(h)

	return w, h
}

func (s *scaler) addParent(x, y ContentAdder) {
	s.xp = x
	s.yp = y
	s.x = s.xc.create(x)
	s.y = s.yc.create(y)
}

func newScaler(x, y createScaler) *scaler {
	ret := &scaler{
		xc: x,
		yc: y,
		xs: &baseSizer{},
		ys: &baseSizer{},
	}
	ret.addParent(nil, nil)

	return ret
}

func newLinearScaler(x, y createScaler) *scaler {
	ret := &scaler{
		xc: x,
		yc: y,
		xs: &parentSizer{},
		ys: &baseSizer{},
	}
	ret.addParent(nil, nil)

	return ret
}

func newListScaler(x, y createScaler) *scaler {
	ret := &scaler{
		xc: x,
		yc: y,
		xs: &baseSizer{},
		ys: &parentSizer{},
	}
	ret.addParent(nil, nil)

	return ret
}
