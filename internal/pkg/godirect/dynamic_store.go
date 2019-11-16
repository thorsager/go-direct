package godirect

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DynamicDirect struct {
	id         uint64
	statusCode int
	path       string
	targetUrl  string
}

func (d *DynamicDirect) URL() string {
	return d.targetUrl
}

func (d *DynamicDirect) Code() int {
	return d.statusCode
}

func (d *DynamicDirect) String() string {
	return fmt.Sprintf("%s(%d) -> (%d)%s", d.path, d.id, d.statusCode, d.targetUrl)
}

func (d *DynamicDirect) Path() string {
	return d.path
}

type DynamicDirectStore struct {
	path2urlFunc path2urlFunc
	store        map[uint64]DynamicDirect
}

func NewDynamicDirectStore() *DynamicDirectStore {
	return &DynamicDirectStore{}
}

func trimPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}

func (s *DynamicDirectStore) Lookup(path string) (Direct, error) {
	if s.store == nil {
		return nil, &NotFoundError{path: path}
	}
	tp := trimPath(path)
	id, err := asNumber(tp)
	if err != nil {
		return nil, &NotFoundError{path: path}
	}
	if d, ok := s.store[id]; ok {
		return &d, nil
	}
	return nil, &NotFoundError{path: path}
}

func (s *DynamicDirectStore) All() []Direct {
	var all []Direct
	for _, value := range s.store {
		all = append(all, &value)
	}
	return all
}

func (s *DynamicDirectStore) Path2UrlFunc(urlFunc path2urlFunc) {
	s.path2urlFunc = urlFunc
}

//func (s *DynamicDirectStore) HandleFunc(w http.ResponseWriter, r *http.Request) {
//	switch r.Method {
//	case http.MethodPost:
//		s.createDirect(w, r)
//	default:
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//	}
//}

func (s *DynamicDirectStore) newId() (uint64, error) {
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

func (s *DynamicDirectStore) CreateAndAdd(code int, targetUrl *url.URL) (Direct, error) {
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

func (s *DynamicDirectStore) Create(code int, targetUrl *url.URL) (Direct, error) {
	id, err := s.newId()
	if err != nil {
		return nil, err
	}
	direct := &DynamicDirect{
		id:         id,
		statusCode: code,
		targetUrl:  targetUrl.String(),
		path:       asString(id),
	}
	return direct, nil
}

func (s *DynamicDirectStore) Add(direct Direct) error {
	dd, ok := direct.(*DynamicDirect)
	if !ok {
		return fmt.Errorf("unable to do casting to *DynamicDirect")
	}
	if s.store == nil {
		s.store = make(map[uint64]DynamicDirect)
	}
	if _, found := s.store[dd.id]; !found {
		s.store[dd.id] = *dd
	} else {
		return fmt.Errorf("storage id clash (%d)", dd.id)
	}
	return nil
}

func (s *DynamicDirectStore) Remove(path string) error {
	d, err := s.Lookup(path)
	if err != nil {
		return err
	}
	dd, ok := d.(*DynamicDirect)
	if !ok {
		return fmt.Errorf("unable to do casting to *DynamicDirect")
	}
	delete(s.store, dd.id)
	return nil
}

func in(needle interface{}, haystack ...interface{}) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}

func asString(number uint64) string {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, number)
	return base64.RawURLEncoding.EncodeToString(buffer)
}

func asNumber(uint64base string) (uint64, error) {
	if len(uint64base) != 11 {
		return 0, fmt.Errorf("invalid string length")
	}
	buffer, err := base64.RawURLEncoding.DecodeString(uint64base)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buffer), nil
}
