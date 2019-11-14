package godirect

import (
	"encoding/json"
	"fmt"
)

type ExactMatchStore struct {
	path2urlFunc path2urlFunc
	redirectMap  map[string]string
	code         int
}

type ExactMatchDirect struct {
	code int
	url  string
	path string
}

func (d *ExactMatchDirect) URL() string {
	return d.url
}

func (d *ExactMatchDirect) Code() int {
	return d.code
}

func (d *ExactMatchDirect) String() string {
	return fmt.Sprintf("%s -> (%d)%s", d.path, d.code, d.url)
}

func ExactMatchStoreFromJson(code int, jsonString string) (*ExactMatchStore, error) {
	redirects := make(map[string]string)
	err := json.Unmarshal([]byte(jsonString), &redirects)
	if err != nil {
		return nil, err
	}
	return &ExactMatchStore{code: code, redirectMap: redirects}, nil
}

func (s *ExactMatchStore) Lookup(path string) (Direct, error) {
	if target, ok := s.redirectMap[path]; ok {
		return &ExactMatchDirect{path: path, code: s.code, url: target}, nil
	}
	return nil, &NotFoundError{path: path}
}

func (s *ExactMatchStore) Path2UrlFunc(urlFunc path2urlFunc) {
	s.path2urlFunc = urlFunc
}

func (s *ExactMatchStore) All() []Direct {
	var all []Direct
	for path, url := range s.redirectMap {
		p := path
		if s.path2urlFunc != nil {
			p = s.path2urlFunc(p)
		}
		all = append(all, &ExactMatchDirect{path: p, code: s.code, url: url})
	}
	return all
}
