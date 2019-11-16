package godirect

type path2urlFunc func(path string) string

type Direct interface {
	URL() string
	Code() int
	Path() string
	String() string
}

type DirectStore interface {
	Lookup(path string) (Direct, error)
	All() []Direct
}
