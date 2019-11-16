package main

import (
	"github.com/joho/godotenv"
	"github.com/thorsager/go-direct/internal/pkg/ephemeralDirectStore"
	"github.com/thorsager/go-direct/internal/pkg/godirect"
	"github.com/thorsager/go-direct/internal/pkg/recursiveDirectStore"
	"github.com/thorsager/go-direct/internal/pkg/staticDirectStore"
	"github.com/thorsager/go-direct/internal/pkg/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	App       = kingpin.New("godirectd", "Go Direct").Version(version.GetVersion())
	Debug     = App.Flag("debug", "Debug mode.").Bool()
	Dynamic   = App.Flag("dynamic", "Enable dynamic mode").Short('d').Default("false").Bool()
	rootStore = recursiveDirectStore.New()
	srv       godirect.Server
	port      int
)

const (
	envPermRedir = "REDIRECTS" // to preserve backwards compat
	envTempRedir = "TEMPORARY_REDIRECTS"
	envWebDir    = "WEB_DIR"
	webDir       = "/web/static"
)

func main() {
	kingpin.MustParse(App.Parse(os.Args[1:]))

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading .env file: %v", err)
	}

	bindAddr := os.Getenv("SERVER_IP")
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %s", portStr)
	}

	log.Printf("Starting go-direct(build:%s)", version.GetVersion())
	srv = godirect.New(bindAddr, port, rootStore)

	if *Dynamic {
		setupDynamic()
	}
	setupStatic()

	if *Debug {
		for _, d := range rootStore.All() {
			log.Printf("DEBUG: %s", d)
		}
	}
	srv.Router().Use(loggingMiddleware)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Unable to start godirect: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Do stuff here
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		log.Printf("Served [%v] %s '%s', in %v", r.RemoteAddr, r.Method, r.RequestURI, time.Since(startTime))
	})
}

func setupStatic() {
	if json := os.Getenv(envPermRedir); json != "" {
		store, err := staticDirectStore.FromJson(http.StatusMovedPermanently, json)
		if err != nil {
			log.Fatalf("Unable to create DirectStore: %v", err)
		}
		rootStore.Add(store)
	}
	if json := os.Getenv(envTempRedir); json != "" {
		store, err := staticDirectStore.FromJson(http.StatusTemporaryRedirect, json)
		if err != nil {
			log.Fatalf("Unable to create DirectStore: %v", err)
		}
		rootStore.Add(store)
	}
}

func setupDynamic() {
	directorUrlStr := os.Getenv("DIRECTOR_URL")
	if directorUrlStr == "" {
		log.Fatal("Director url is required, please set DIRECTOR_URL")
	}

	directorUrl, err := url.Parse(directorUrlStr)
	if err != nil {
		log.Fatalf("invalid DIRECTOR_URL: %s", err)
	}

	siteUrlStr := os.Getenv("SITE_URL")
	if siteUrlStr == "" {
		log.Fatal("Site url is required, please set SITE_URL")
	}

	siteUrl, err := url.Parse(siteUrlStr)
	if err != nil {
		log.Fatalf("invalid DIRECTOR_URL: %s", err)
	}

	dStore := ephemeralDirectStore.New()

	rootStore.Add(dStore)

	dir := webDir
	if d := os.Getenv(envWebDir); d != "" {
		dir = d
	}

	srv.Router().Host(siteUrl.Hostname()).PathPrefix("/api").HandlerFunc(godirect.DynamicDirectHandlerFunc(directorUrl, dStore))
	srv.Router().Host(siteUrl.Hostname()).PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	log.Printf("  directorURL: %s", directorUrl)
	log.Printf("      siteURL: %s", siteUrl)
}
