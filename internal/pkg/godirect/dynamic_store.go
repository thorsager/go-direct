package godirect

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DynamicDirect struct {
	Id         uint64 `json:"id"`
	StatusCode int    `json:"code"`
	Url        string `json:"url"`
	Path       string `json:"path"`
}

func (d *DynamicDirect) URL() string {
	return d.Url
}

func (d *DynamicDirect) Code() int {
	return d.StatusCode
}

func (d *DynamicDirect) String() string {
	return fmt.Sprintf("%s(%d) -> (%d)%s", d.Path, d.Id, d.StatusCode, d.Url)
}

type DynamicDirectStore struct {
	path2urlFunc path2urlFunc
	store        map[uint64]DynamicDirect
}

func NewDynamicDirectStore() *DynamicDirectStore {
	return &DynamicDirectStore{}
}

func (s *DynamicDirectStore) Lookup(path string) (Direct, error) {
	if s.store == nil {
		return nil, &NotFoundError{path: path}
	}

	id, err := asNumber(path[1:])
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

type CreateRequest struct {
	Code int    `json:"code"`
	Url  string `json:"url"`
}

func (s *DynamicDirectStore) HandleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createDirect(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

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

func (s *DynamicDirectStore) createDirect(w http.ResponseWriter, r *http.Request) {
	cRequest := &CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(cRequest)
	if err != nil {
		http.Error(w, "Bad Request\n"+err.Error(), http.StatusBadRequest)
		return
	}
	targetUrl, err := url.Parse(cRequest.Url)
	if err != nil {
		http.Error(w, "Bad Request\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if !in(targetUrl.Scheme, "http", "https") {
		http.Error(w, "Bad Request\nInvalid url scheme (http,https)", http.StatusBadRequest)
		return
	}

	if !in(cRequest.Code, http.StatusTemporaryRedirect, http.StatusMovedPermanently) {
		http.Error(w, "Bad Request\nInvalid code (301,307)", http.StatusBadRequest)
		return
	}

	id, err := s.newId()
	if err != nil {
		http.Error(w, "Internal Server Error\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	p := asString(id)
	if s.path2urlFunc != nil {
		p = s.path2urlFunc(p)
	}
	d := &DynamicDirect{
		Id:         id,
		StatusCode: cRequest.Code,
		Url:        targetUrl.String(),
		Path:       s.path2urlFunc(asString(id)),
	}
	if s.store == nil {
		s.store = make(map[uint64]DynamicDirect)
	}
	s.store[d.Id] = *d
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		http.Error(w, "Internal Server Error\n"+err.Error(), http.StatusInternalServerError)
	}
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
	buffer, err := base64.RawURLEncoding.DecodeString(uint64base)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buffer), nil
}
