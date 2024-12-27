package chocolate

type ContentSizer interface {
	AddContentX(width int)
	AddContentY(height int)
}

type ParentSizer interface {
	AddPartsX(parts int)
	AddPartsY(parts int)
	TakePartsX(parts int) int
	TakePartsY(parts int) int
}

type Scaler interface {
	finalSize(size int) int
	setCotentSize(size int) int
}

type fixed struct {
	value         int
	addContentFct func(int)
}

func (s *fixed) finalSize(size int) int { return s.value }
func (s *fixed) setCotentSize(size int) int {
	if s.addContentFct != nil {
		s.addContentFct(s.value)
	}
	return s.value
}

type dynamic struct {
	addContentFct func(int)
}

func (s *dynamic) finalSize(size int) int { return size }
func (s *dynamic) setCotentSize(size int) int {
	if s.addContentFct != nil {
		s.addContentFct(size)
	}
	return size
}

type parent struct {
	value         int
	addContentFct func(int)
	take          func(int) int
	add           func(int)
}

func (s *parent) finalSize(size int) int {
	return s.take(s.value)
}

func (s *parent) setCotentSize(size int) int {
	if s.addContentFct != nil {
		s.addContentFct(size)
	}
	return 0
}

type scaler struct {
	x Scaler
	y Scaler
	p ContentSizer
}

func (s *scaler) FinalSize(width, height int) (int, int) {
	w := s.x.finalSize(width)
	h := s.y.finalSize(height)

	return w, h
}

func (s *scaler) SetContentSize(width, height int) (int, int) {
	w, h := width, height
	if s.p != nil {
		w = s.x.setCotentSize(width)
		h = s.y.setCotentSize(height)
	}

	return w, h
}

type scalerOption func(*scaler)

func withXfixed(v int) scalerOption {
	return func(s *scaler) {
		sp := &fixed{
			value: v,
		}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentX
		}
		s.x = sp
	}
}

func withYfixed(v int) scalerOption {
	return func(s *scaler) {
		sp := &fixed{
			value: v,
		}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentY
		}
		s.y = sp
	}
}

func withXdynamic() scalerOption {
	return func(s *scaler) {
		sp := &dynamic{}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentX
		}
		s.x = sp
	}
}

func withYdynamic() scalerOption {
	return func(s *scaler) {
		sp := &dynamic{}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentY
		}
		s.y = sp
	}
}

func withXparent(v int, p ParentSizer) scalerOption {
	return func(s *scaler) {
		sp := &parent{
			value: v,
		}
		if p != nil {
			sp.take = p.TakePartsX
			sp.add = p.AddPartsX
		}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentX
		}
		s.x = sp
	}
}

func withYparent(v int, p ParentSizer) scalerOption {
	return func(s *scaler) {
		sp := &parent{
			value: v,
		}
		if p != nil {
			sp.take = p.TakePartsY
			sp.add = p.AddPartsY
		}
		if s.p != nil {
			sp.addContentFct = s.p.AddContentY
		}
		s.y = sp
	}
}

func newScaler(parent ContentSizer, opts ...scalerOption) *scaler {
	ret := &scaler{
		p: parent,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret
}
