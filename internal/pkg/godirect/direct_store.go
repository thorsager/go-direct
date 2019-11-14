package godirect

type path2urlFunc func(path string) string

type Direct interface {
	URL() string
	Code() int
	String() string
}

type DirectStore interface {
	Lookup(path string) (Direct, error)
	All() []Direct
	Path2UrlFunc(p2uFunc path2urlFunc)
}
