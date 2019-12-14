package fileDirectStore

import (
	"fmt"
	"github.com/thorsager/go-direct/internal/pkg/godirect"
	"github.com/thorsager/go-direct/internal/pkg/structFileStore"
	"github.com/thorsager/go-direct/internal/pkg/utl"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type fileDirect struct {
	ForwardCode int    `json:"status_code"`
	DirectPath  string `json:"path"`
	TargetUrl   string `json:"target_url"`
	PathId      uint64 `json:"id"`
}

func (d *fileDirect) Id() interface{} {
	return d.PathId
}

func (d *fileDirect) URL() string {
	return d.TargetUrl
}

func (d *fileDirect) Code() int {
	return d.ForwardCode
}

func (d *fileDirect) String() string {
	return fmt.Sprintf("%s(%d) -> (%d)%s", d.DirectPath, d.PathId, d.ForwardCode, d.TargetUrl)
}

func (d *fileDirect) Path() string {
	return d.DirectPath
}

type FileDirectStore struct {
	structStore *structFileStore.StructFileStore
}

func New(location string) (*FileDirectStore, error) {
	structStore, err := structFileStore.NewJSON(location)
	if err != nil {
		return nil, err
	}
	return &FileDirectStore{structStore: structStore}, nil
}

func (s *FileDirectStore) Lookup(path string) (godirect.Direct, error) {
	tp := utl.TrimPath(path)
	d := &fileDirect{}
	err := s.structStore.Get(tp, d)
	if err != nil {
		return nil, godirect.NotFound(path)
	}
	return d, nil
}

func (s *FileDirectStore) All() ([]godirect.Direct, error) {
	var all []godirect.Direct
	ids, err := s.structStore.All()
	if err != nil {
		return all, err
	}
	defer ids.Close()

	for ids.Next() {
		var fd fileDirect
		err := ids.Scan(&fd)
		if err != nil {
			return all, err
		}
		all = append(all, &fd)
	}
	return all, nil
}

func (s *FileDirectStore) newId() (uint64, error) {
	maxCount := 10
	for maxCount > 0 {
		id := rand.Uint64()
		if !s.structStore.Exist(utl.AsBase64String(id)) {
			return id, nil
		}
		maxCount = maxCount - 1
	}
	return 0, fmt.Errorf("unable to find availabel id in %d rounds", maxCount)
}

func (s *FileDirectStore) CreateAndAdd(code int, targetUrl *url.URL) (godirect.IdentifiedDirect, error) {
	d, err := s.Create(code, targetUrl)
	if err != nil {
		return nil, err
	}
	err = s.Add(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (s *FileDirectStore) Create(code int, targetUrl *url.URL) (godirect.IdentifiedDirect, error) {
	id, err := s.newId()
	if err != nil {
		return nil, err
	}
	direct := &fileDirect{
		ForwardCode: code,
		TargetUrl:   targetUrl.String(),
		DirectPath:  utl.AsBase64String(id),
		PathId:      id,
	}
	return direct, nil
}

func (s *FileDirectStore) Add(direct godirect.IdentifiedDirect) error {
	if s.structStore == nil {
		return fmt.Errorf("store not initialized")
	}
	ed, ok := direct.(*fileDirect)
	if !ok {
		return fmt.Errorf("unable to do casting to *ephemeralDirect")
	}
	err := s.structStore.Put(ed.Path(), ed)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileDirectStore) Remove(path string) error {
	if !s.structStore.Exist(path) {
		return godirect.NotFound(path)
	}
	err := s.structStore.Delete(path)
	if err != nil {
		return err
	}
	return nil
}
