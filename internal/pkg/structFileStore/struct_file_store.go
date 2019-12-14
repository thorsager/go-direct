package structFileStore

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path"
)

const perm = 0600
const dirPerm = 0700

type MarshalFunc func(o interface{}) ([]byte, error)
type UnmarshalFunc func(data []byte, o interface{}) error

type StructFileStore struct {
	baseLocation string
	extension    string
	marshal      MarshalFunc
	unMarshal    UnmarshalFunc
}

func New(path string, marshalFunc MarshalFunc, unmarshalFunc UnmarshalFunc, extension string) (*StructFileStore, error) {
	err := os.MkdirAll(path, dirPerm)
	if err != nil {
		return nil, err
	}
	return &StructFileStore{
		baseLocation: path,
		marshal:      marshalFunc,
		unMarshal:    unmarshalFunc,
		extension:    extension,
	}, nil

}

func NewJSON(path string) (*StructFileStore, error) {
	return New(path, json.Marshal, json.Unmarshal, ".json")
}

func (f *StructFileStore) Delete(id interface{}) error {
	itemPath := f.asPath(asKey(id))
	if !exists(itemPath) {
		return fmt.Errorf("item not found (%v): %s", id, itemPath)
	}
	return os.Remove(itemPath)
}

func (f *StructFileStore) All() (*StoreIterator, error) {
	return newStoreIterator(f)
}

func (f *StructFileStore) Put(id interface{}, o interface{}) error {
	itemPath := f.asPath(asKey(id))
	if exists(itemPath) {
		return fmt.Errorf("item already exist (%v): %s", id, itemPath)
	}
	fileContent, err := f.marshal(o)
	if err != nil {
		return err
	}
	err = ensurePath(itemPath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(itemPath, fileContent, perm)
}

func (f *StructFileStore) Get(id interface{}, o interface{}) error {
	itemPath := f.asPath(asKey(id))
	return f.getByItemPath(itemPath, o)
}

func (f *StructFileStore) Exist(id interface{}) bool {
	itemPath := f.asPath(asKey(id))
	return exists(itemPath)
}

func (f *StructFileStore) getByItemPath(itemPath string, o interface{}) error {
	fileContent, err := ioutil.ReadFile(itemPath)
	if err != nil {
		return err
	}
	return f.unMarshal(fileContent, o)
}

func exists(itemPath string) bool {
	_, err := os.Stat(itemPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func asKey(i interface{}) string {
	raw := fmt.Sprintf("%v", i)
	h := fnv.New64a()
	_, _ = h.Write([]byte(raw))
	return fmt.Sprintf("%.16x", h.Sum64())
}

func ensurePath(itemPath string) error {
	dir, _ := path.Split(itemPath)
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return err
	}
	return nil
}

func (f *StructFileStore) asPath(itemName string) string {
	return path.Join(f.baseLocation, itemName) + f.extension
}
