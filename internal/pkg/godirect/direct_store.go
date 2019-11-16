package godirect

import "net/url"

type Direct interface {
	URL() string
	Code() int
	Path() string
	String() string
}

type PersistableDirect interface {
	Direct
	Id() interface{}
}

type DirectStore interface {
	Lookup(path string) (Direct, error)
	All() []Direct
}

type MutableDirectStore interface {
	DirectStore
	Remove(path string) error
	Add(direct PersistableDirect) error
	Create(code int, targetUrl *url.URL) (PersistableDirect, error)
	CreateAndAdd(code int, targetUrl *url.URL) (PersistableDirect, error)
}
