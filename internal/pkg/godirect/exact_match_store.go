package godirect

import (
	"encoding/json"
	"fmt"
)

type ExactMatchStore struct {
	redirectMap map[string]string
	code        int
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

func (d *ExactMatchDirect) Path() string {
	return d.path
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

func (s *ExactMatchStore) All() []Direct {
	var all []Direct
	for path, url := range s.redirectMap {
		all = append(all, &ExactMatchDirect{path: path, code: s.code, url: url})
	}
	return all
}
