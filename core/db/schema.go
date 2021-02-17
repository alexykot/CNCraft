// Code generated for package db by go-bindata DO NOT EDIT. (@generated)
// sources:
// schema/001_cncraft.down.sql
// schema/001_players.up.sql
package db

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __001_cncraftDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x76\xf6\x70\xf5\x75\x54\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xce\x4b\x2e\x4a\x4c\x2b\x51\x70\x76\x0c\x76\x76\x74\x71\xb5\xe6\x02\x2b\x0b\x71\x74\xf2\x71\x45\x52\x55\x50\x9a\x94\x93\x99\xac\x07\x55\x1c\x5f\x9c\x9c\x91\x9a\x9b\x18\x9f\x9b\x99\x5e\x94\x58\x92\x99\x9f\x57\x6c\xcd\x05\x08\x00\x00\xff\xff\xd6\x2c\x99\x35\x5e\x00\x00\x00")

func _001_cncraftDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__001_cncraftDownSql,
		"001_cncraft.down.sql",
	)
}

func _001_cncraftDownSql() (*asset, error) {
	bytes, err := _001_cncraftDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "001_cncraft.down.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xc1\x4f\xc2\x30\x18\xc5\xef\xfd\x2b\x5e\x38\x41\x22\x26\x7a\x32\xe1\x54\x47\x09\x8d\xdb\x4a\xba\x4e\xc1\x0b\xa9\x5b\x71\x8b\xb0\x2d\x6d\x17\x9d\x7f\xbd\x61\x3a\x0f\x26\x84\xe0\x3b\xb6\xef\xf7\xde\xf7\x7d\xd3\x29\x54\x61\xe0\xb2\xc2\x1c\x34\x76\xb5\xc5\x8b\xae\xde\xb0\x37\xf9\xab\xb1\xc8\xb5\xd7\x24\x90\x8c\x2a\x86\x24\x58\xb2\x88\x82\x2f\x10\x0b\x05\xb6\xe6\x89\x4a\x90\x55\x99\xd5\x3b\x3f\x1b\x4c\x6c\xad\x58\x9c\x70\x11\xff\xf1\x8d\xda\xb6\xcc\xa7\xb5\x73\xcd\x68\x46\x06\xb3\xa2\xf7\x21\x1b\x22\xae\x9b\xbd\xee\x8c\x75\x64\x4c\x00\xa0\xcc\x31\x28\x4d\xf9\x1c\x27\x74\xac\x88\xd3\x30\xc4\x4a\xf2\x88\xca\x0d\x1e\xd8\xe6\xaa\x0f\x68\x9d\xb1\x95\x3e\x18\x00\x8f\x54\x06\x4b\x2a\xc7\x37\xb7\x77\x93\x53\x01\xdf\x50\x53\xbb\xd2\x97\x75\xb5\xfd\xc0\x5c\xa4\xc7\xe9\x56\x92\x05\xbc\x5f\xe8\x3c\xd4\xfd\x07\xfa\xbc\x00\xea\xf4\xfb\xf0\xb1\x08\x05\x55\x67\x8e\xf2\xd3\x54\xfa\xac\xb8\x14\xca\xac\xd1\xde\xe4\x5b\xed\xa1\x78\xc4\x12\x45\xa3\x15\x9e\xb8\x5a\x8a\x54\xf5\x2f\x78\x16\x31\xfb\x85\xc8\x64\x46\xbe\x02\x00\x00\xff\xff\x41\x69\x57\x3f\x4a\x02\x00\x00")

func _001_playersUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__001_playersUpSql,
		"001_players.up.sql",
	)
}

func _001_playersUpSql() (*asset, error) {
	bytes, err := _001_playersUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "001_players.up.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"001_cncraft.down.sql": _001_cncraftDownSql,
	"001_players.up.sql":   _001_playersUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"001_cncraft.down.sql": &bintree{_001_cncraftDownSql, map[string]*bintree{}},
	"001_players.up.sql":   &bintree{_001_playersUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
