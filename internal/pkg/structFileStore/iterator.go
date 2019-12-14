package structFileStore

import (
	"io"
	"os"
	"path"
)

//type FilterFunc func(item interface{}) bool
type StoreIterator struct {
	store *StructFileStore
	dh    *os.File

	cInfo os.FileInfo
	cErr  error
}

func newStoreIterator(store *StructFileStore) (*StoreIterator, error) {
	i := &StoreIterator{store: store}

	fh, err := os.Open(i.store.baseLocation)
	if err != nil {
		return nil, err
	}
	i.dh = fh

	return i, nil
}

func (i *StoreIterator) readNext() {
	dirInfoList, err := i.dh.Readdir(1)
	if err != nil {
		i.cErr = err
		i.cInfo = nil
	} else {
		i.cInfo = dirInfoList[0]
	}
}

func (i *StoreIterator) Close() error {
	if i.dh != nil {
		return i.dh.Close()
	}
	return nil
}

func (i *StoreIterator) Next() bool {
	i.readNext()
	return i.cErr != io.EOF
}

func (i *StoreIterator) Scan(o interface{}) error {
	if i.cErr != nil {
		return i.cErr
	}
	itemPath := path.Join(i.store.baseLocation, i.cInfo.Name())
	return i.store.getByItemPath(itemPath, o)
}
