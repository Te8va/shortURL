package repository

type MapStore struct {
	data map[string]string
}

func NewMapStore() *MapStore {
	return &MapStore{
		data: make(map[string]string),
	}
}

func (s *MapStore) Save(id, url string) error {
	s.data[id] = url
	return nil
}

func (s *MapStore) Get(id string) (string, bool) {
	url, exists := s.data[id]
	return url, exists
}
