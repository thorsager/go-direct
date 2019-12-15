package logging

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func (rw *statusRecorder) status() string {
	return fmt.Sprintf("%d %s", rw.statusCode, http.StatusText(rw.statusCode))
}
func (rw *statusRecorder) is(statusCode int) bool {
	return rw.statusCode == statusCode
}

func wrapWriter(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{w, http.StatusOK}
}

func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lw := wrapWriter(w)
		next.ServeHTTP(lw, r)
		if lw.is(http.StatusTemporaryRedirect) || lw.is(http.StatusMovedPermanently) || lw.is(http.StatusPermanentRedirect) {
			location := w.Header().Get("location")
			log.Printf("[%v] %s '%s' -> '%s' [%s], in %v", r.RemoteAddr, r.Method, r.RequestURI, location, lw.status(), time.Since(startTime))
		} else {
			log.Printf("[%v] %s '%s' [%s], in %v", r.RemoteAddr, r.Method, r.RequestURI, lw.status(), time.Since(startTime))
		}
	})
}
