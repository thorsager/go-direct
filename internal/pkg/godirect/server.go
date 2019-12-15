package godirect

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	srv      http.Server
	router   *mux.Router
	store    DirectStore
	hostname *string
}

func New(bindAddr string, port int, store DirectStore) Server {
	router := mux.NewRouter()
	router.StrictSlash(true)
	s := Server{router: router, store: store}

	s.srv = http.Server{
		Addr:         fmt.Sprintf("%s:%d", bindAddr, port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return s
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) Hostname(name string) {
	s.hostname = &name
}

func (s *Server) addFallbackHandler() {
	if s.hostname != nil {
		s.router.Host(*s.hostname).PathPrefix("/").HandlerFunc(s.redirect)
	} else {
		s.router.PathPrefix("/").HandlerFunc(s.redirect)
	}
}

func (s *Server) ListenAndServe() error {
	s.addFallbackHandler()
	return s.srv.ListenAndServe()
}
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	s.addFallbackHandler()
	return s.srv.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if target, err := s.store.Lookup(path); err == nil {
		http.Redirect(w, r, target.URL(), target.Code())
	} else {
		http.NotFound(w, r)
	}
}
