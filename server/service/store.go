package service

type fileMeta struct {
	ids  []string
	size int
}

type Store map[string]fileMeta

func NewStore() Store {
	return make(map[string]fileMeta)
}

func (s Store) SetPath(path string, size int, ids []string) {
	s[path] = fileMeta{size: size, ids: ids}
}

func (s Store) GetPath(path string) (int, []string) {
	fm := s[path]
	return fm.size, fm.ids
}
