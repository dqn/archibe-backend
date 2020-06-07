package lib

type Memo struct {
	m    map[string]struct{}
	size uint64
}

func NewMemo(size uint64) *Memo {
	return &Memo{make(map[string]struct{}, size), size}
}

func (s *Memo) Add(key string) {
	s.m[key] = struct{}{}
}

func (s *Memo) Exists(key string) bool {
	_, ok := s.m[key]
	return ok
}

func (s *Memo) Clear(key string) {
	s.m = make(map[string]struct{}, s.size)
}
