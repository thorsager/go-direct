package recursiveDirectStore

import "github.com/thorsager/go-direct/internal/pkg/godirect"

type RecursiveDirectStore struct {
	stores []godirect.DirectStore
}

func New(stores ...godirect.DirectStore) *RecursiveDirectStore {
	return &RecursiveDirectStore{stores: stores}
}

func (m *RecursiveDirectStore) Lookup(path string) (godirect.Direct, error) {
	for _, s := range m.stores {
		target, err := s.Lookup(path)
		if err != nil {
			if godirect.IsNotFound(err) {
				continue
			}
			return nil, err
		}
		return target, nil
	}
	return nil, godirect.NotFound(path)
}

func (m *RecursiveDirectStore) All() ([]godirect.Direct, error) {
	var all []godirect.Direct
	for _, store := range m.stores {
		items, err := store.All()
		if err != nil {
			return all, err
		}
		all = append(all, items...)
	}
	return all, nil
}

func (m *RecursiveDirectStore) Add(store godirect.DirectStore) {
	m.stores = append(m.stores, store)
}

func (m *RecursiveDirectStore) Remove(store godirect.DirectStore) bool {
	if i, found := m.index(store); found {
		m.stores = append(m.stores[:i], m.stores[i+1:]...)
	}
	return false
}

func (m *RecursiveDirectStore) index(store godirect.DirectStore) (int, bool) {
	for idx, s := range m.stores {
		if s == store {
			return idx, true
		}
	}
	return -1, false
}
