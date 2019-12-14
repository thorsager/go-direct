package staticDirectStore

import (
	"encoding/json"
	"fmt"
	"github.com/thorsager/go-direct/internal/pkg/godirect"
)

type StaticDirectStore struct {
	redirectMap map[string]string
	code        int
}

type simpleDirect struct {
	statusCode int
	targetUrl  string
	path       string
}

func (d *simpleDirect) URL() string {
	return d.targetUrl
}

func (d *simpleDirect) Code() int {
	return d.statusCode
}

func (d *simpleDirect) Path() string {
	return d.path
}

func (d *simpleDirect) String() string {
	return fmt.Sprintf("%s -> (%d)%s", d.path, d.statusCode, d.targetUrl)
}

func FromJson(code int, jsonString string) (*StaticDirectStore, error) {
	redirects := make(map[string]string)
	err := json.Unmarshal([]byte(jsonString), &redirects)
	if err != nil {
		return nil, err
	}
	return &StaticDirectStore{code: code, redirectMap: redirects}, nil
}

func (s *StaticDirectStore) Lookup(path string) (godirect.Direct, error) {
	if target, ok := s.redirectMap[path]; ok {
		return &simpleDirect{path: path, statusCode: s.code, targetUrl: target}, nil
	}
	return nil, godirect.NotFound(path)
}

func (s *StaticDirectStore) All() ([]godirect.Direct, error) {
	var all []godirect.Direct
	for path, url := range s.redirectMap {
		all = append(all, &simpleDirect{path: path, statusCode: s.code, targetUrl: url})
	}
	return all, nil
}
