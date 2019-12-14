package structFileStore

type FileStore interface {
	Put(interface{}, interface{}) error
	Get(interface{}, interface{}) error
	Delete(interface{}) error
	Iterator()
}
