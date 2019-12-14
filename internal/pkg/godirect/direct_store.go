package godirect

import (
	"net/url"
)

type Direct interface {
	URL() string
	Code() int
	Path() string
	String() string
}

type IdentifiedDirect interface {
	Direct
	Id() interface{}
}

type DirectStore interface {
	Lookup(path string) (Direct, error)
	All() ([]Direct, error)
}

type MutableDirectStore interface {
	DirectStore
	Remove(path string) error
	Add(direct IdentifiedDirect) error
	Create(code int, targetUrl *url.URL) (IdentifiedDirect, error)
	CreateAndAdd(code int, targetUrl *url.URL) (IdentifiedDirect, error)
}
