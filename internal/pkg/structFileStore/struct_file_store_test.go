package structFileStore

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"reflect"
	"testing"
)

type Data struct {
	Foo int    `json:"foo"`
	Bar string `json:"bar"`
}

type NJData struct {
	Bar string
}

func setup() (*StructFileStore, string, func()) {
	tmpDir, _ := ioutil.TempDir("", "_test")
	store, err := NewJSON(tmpDir)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return store, tmpDir, func() { _ = os.RemoveAll(tmpDir) }
}

func createAndAddRandom(store *StructFileStore) (interface{}, string, *Data) {
	key := rand.Int()
	d := &Data{Foo: key, Bar: fmt.Sprintf("Hello World! (%d)", key)}
	err := store.Put(key, d)
	if err != nil {
		log.Fatalf("unable to create testdata: %v", err)
	}
	return key, store.asPath(asKey(key)), d
}

func cat(file string) string {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read content of %v", file)
	}
	return string(body)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestStructFileStore_All(t *testing.T) {
	f, _, cleanup := setup()
	defer cleanup()

	for i := 0; i < 100; i++ {
		_, _, _ = createAndAddRandom(f)
	}

	it, err := f.All()
	if err != nil {
		t.Errorf("unable to create iterator: %v", err)
	}
	defer it.Close()

	for it.Next() {
		var item Data
		err := it.Scan(&item)
		if err != nil {
			t.Errorf("unable to scan item: %v", err)
		}
		fmt.Printf("%+v\n", item)
	}

}

func TestStructFileStore_Delete(t *testing.T) {
	f, _, cleanup := setup()
	defer cleanup()

	key, file, _ := createAndAddRandom(f)
	if !exists(file) {
		t.Errorf("initial store-file missign: %s", file)
	}

	err := f.Delete(key)
	if err != nil {
		t.Errorf("unable to delte item=%v: %v", key, err)
	}

	if exists(file) {
		t.Errorf("file for item=%v still on file-sytem: %s", key, file)
	}
}

func TestStructFileStore_Put(t *testing.T) {
	f, _, cleanup := setup()
	defer cleanup()

	key := rand.Int()

	d := &Data{Foo: 0, Bar: "Hello World!"}
	err := f.Put(key, d)
	if err != nil {
		t.Errorf("failed to Put '%+v' to id=%v, %v", d, key, err)
	}

	storeFile := f.asPath(asKey(key))
	if !fileExists(storeFile) {
		t.Errorf("failed locate storage file %s for id=%v", storeFile, key)
	}

	fileContent := cat(storeFile)
	expected, _ := json.Marshal(d)
	if fileContent != string(expected) {
		t.Errorf("content mismatch %s vs. %s", fileContent, expected)
	}
}

func TestStructFileStore_Put_NJ(t *testing.T) {
	f, _, cleanup := setup()
	defer cleanup()

	key := rand.Int()

	d := &NJData{Bar: "Hello World!"}
	err := f.Put(key, d)
	if err != nil {
		t.Errorf("failed to Put '%+v' to id=%v, %v", d, key, err)
	}

	storeFile := f.asPath(asKey(key))
	if !fileExists(storeFile) {
		t.Errorf("failed locate storage file %s for id=%v", storeFile, key)
	}
}

func TestStructFileStore_Get(t *testing.T) {
	f, _, cleanup := setup()
	defer cleanup()

	key, file, d := createAndAddRandom(f)
	if !exists(file) {
		t.Errorf("initial store-file missign: %s", file)
	}

	d1 := &Data{}
	d2 := reflect.New(reflect.TypeOf("Data"))
	fmt.Printf("%+v", d2)
	err := f.Get(key, d1)
	if err != nil {
		t.Errorf("failed to get id=%v, %v", key, err)
	}

	if diff := deep.Equal(d, d1); diff != nil {
		t.Error(diff)
	}
}

func TestStructFileStore_asPath(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "_test")
	type fields struct {
		path string
	}
	type args struct {
		itemName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"relative", fields{"rel"}, args{"af63ad4c86019caf"}, "rel/af63ad4c86019caf.json"},
		{"absolute", fields{tmpDir}, args{"af63ad4c86019caf"}, path.Join(tmpDir, "af63ad4c86019caf.json")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewJSON(tt.fields.path)
			if err != nil {
				log.Fatalf("%v", err)
			}
			if got := f.asPath(tt.args.itemName); got != tt.want {
				t.Errorf("asPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_asKey(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"zero", args{0}, "af63ad4c86019caf"},
		{"int", args{2342}, "20f53f0b3a7f8460"},
		{"uint", args{uint64(2142)}, "0f1a7d0b303c45ce"},
		{"string", args{"hello world"}, "779a65e7023cd2e7"},
		{"slice", args{[]int{1, 2, 3}}, "20abce5cba4770b9"},
		{"map", args{map[string]int{"a": 97}}, "160fcb17f5e2059c"},
		{"struct", args{&Data{Foo: 1, Bar: "joe"}}, "74215c6a18a5dd16"},
		{"bool", args{true}, "5b5c98ef514dbfa5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := asKey(tt.args.i); got != tt.want {
				t.Errorf("asKey(%v) = %v, want %v", tt.args.i, got, tt.want)
			}
		})
	}
}
