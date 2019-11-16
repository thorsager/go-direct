package ephemeralDirectStore

import (
	"fmt"
	"github.com/thorsager/go-direct/internal/pkg/godirect"
	"github.com/thorsager/go-direct/internal/pkg/utl"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ephemeralDirect struct {
	statusCode int
	path       string
	targetUrl  string
	id         uint64
}

func (d *ephemeralDirect) Id() interface{} {
	return d.id
}

func (d *ephemeralDirect) URL() string {
	return d.targetUrl
}

func (d *ephemeralDirect) Code() int {
	return d.statusCode
}

func (d *ephemeralDirect) String() string {
	return fmt.Sprintf("%s(%d) -> (%d)%s", d.path, d.id, d.statusCode, d.targetUrl)
}

func (d *ephemeralDirect) Path() string {
	return d.path
}

type EphemeralDirectStore struct {
	store map[uint64]ephemeralDirect
}

func New() *EphemeralDirectStore {
	return &EphemeralDirectStore{}
}

func (s *EphemeralDirectStore) Lookup(path string) (godirect.Direct, error) {
	if s.store == nil {
		return nil, godirect.NotFound(path)
	}
	tp := utl.TrimPath(path)
	id, err := utl.AsUint64(tp)
	if err != nil {
		return nil, godirect.NotFound(path)
	}
	if d, ok := s.store[id]; ok {
		return &d, nil
	}
	return nil, godirect.NotFound(path)
}

func (s *EphemeralDirectStore) All() []godirect.Direct {
	var all []godirect.Direct
	for _, value := range s.store {
		all = append(all, &value)
	}
	return all
}

func (s *EphemeralDirectStore) newId() (uint64, error) {
	maxCount := 10
	for maxCount > 0 {
		id := rand.Uint64()
		if _, found := s.store[id]; !found {
			return id, nil
		}
		maxCount = maxCount - 1
	}
	return 0, fmt.Errorf("unable to find availabel id in %d rounds", maxCount)
}

func (s *EphemeralDirectStore) CreateAndAdd(code int, targetUrl *url.URL) (godirect.PersistableDirect, error) {
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

func (s *EphemeralDirectStore) Create(code int, targetUrl *url.URL) (godirect.PersistableDirect, error) {
	id, err := s.newId()
	if err != nil {
		return nil, err
	}
	direct := &ephemeralDirect{
		statusCode: code,
		targetUrl:  targetUrl.String(),
		path:       utl.AsBase64String(id),
		id:         id,
	}
	return direct, nil
}

func (s *EphemeralDirectStore) Add(direct godirect.PersistableDirect) error {
	ed, ok := direct.(*ephemeralDirect)
	if !ok {
		return fmt.Errorf("unable to do casting to *ephemeralDirect")
	}
	if s.store == nil {
		s.store = make(map[uint64]ephemeralDirect)
	}
	if _, found := s.store[ed.id]; !found {
		s.store[ed.id] = *ed
	} else {
		return fmt.Errorf("storage id clash (%d)", ed.id)
	}
	return nil
}

func (s *EphemeralDirectStore) Remove(path string) error {
	d, err := s.Lookup(path)
	if err != nil {
		return err
	}
	dd, ok := d.(*ephemeralDirect)
	if !ok {
		return fmt.Errorf("unable to do casting to *DynamicDirect")
	}
	delete(s.store, dd.id)
	return nil
}
