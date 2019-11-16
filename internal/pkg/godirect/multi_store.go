package godirect

type MultiStore struct {
	path2urlFunc path2urlFunc
	stores       []DirectStore
}

func NewMultiStore(stores ...DirectStore) *MultiStore {
	return &MultiStore{stores: stores}
}

func (m *MultiStore) Lookup(path string) (Direct, error) {
	for _, s := range m.stores {
		target, err := s.Lookup(path)
		if err != nil {
			if IsNotFound(err) {
				continue
			}
			return nil, err
		}
		return target, nil
	}
	return nil, &NotFoundError{path: path}
}

func (m *MultiStore) All() []Direct {
	var all []Direct
	for _, store := range m.stores {
		all = append(all, store.All()...)
	}
	return all
}

func (m *MultiStore) Add(store DirectStore) {
	m.stores = append(m.stores, store)
}
