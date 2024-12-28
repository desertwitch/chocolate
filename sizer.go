package chocolate

type baseSizer struct {
	size        int
	maxSize     int
	contentSize int
}

func (s *baseSizer) GetSize() int          { return s.size }
func (s *baseSizer) GetMaxSize() int       { return s.maxSize }
func (s *baseSizer) GetContentSize() int   { return s.contentSize }
func (s *baseSizer) GetAvailableSize() int { return s.maxSize - s.contentSize }

func (s *baseSizer) SetSize(size int)        { s.size = size }
func (s *baseSizer) SetMaxSize(size int)     { s.maxSize = size }
func (s *baseSizer) SetContentSize(size int) { s.contentSize = size }

func (s *baseSizer) AddContent(size int) { s.contentSize += size }
func (s *baseSizer) AddParts(parts int)  {}
func (s *baseSizer) TakeParts(parts int) int {
	if s.GetSize() > 0 {
		return s.GetSize()
	}
	return s.GetMaxSize()
}

type parentSizer struct {
	baseSizer
	totalParts int
	takenParts int
}

func (s *parentSizer) AddParts(parts int) { s.totalParts += parts }
func (s *parentSizer) TakeParts(parts int) int {
	if s.totalParts <= 0 {
		s.totalParts = 0
		s.takenParts = 0
		return 0
	}

	partSize := (s.GetAvailableSize() / s.totalParts) * parts
	s.takenParts += parts
	if s.takenParts >= s.totalParts {
		partSize += s.size % s.totalParts
		s.totalParts = 0
		s.takenParts = 0
	}

	return partSize
}

func (s *parentSizer) SetSize(size int) {
	s.baseSizer.SetSize(size)
	s.totalParts = 0
	s.takenParts = 0
}

func (s *parentSizer) AddContent(size int) {
	if size > s.GetContentSize() {
		s.SetContentSize(size)
	}
}
