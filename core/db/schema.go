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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xc1\x4f\xc2\x30\x18\xc5\xef\xfb\x2b\x5e\x38\x41\x22\x26\x7a\x32\xe1\x54\xc7\x47\x68\xdc\x5a\xdc\x5a\x05\x2f\x4b\xdd\x8a\x2c\xc2\x46\xb6\x11\xc1\xbf\xde\x30\x29\x26\x24\x84\xe0\x77\xec\xf7\x7e\xef\xb5\x7d\xfd\x3e\xd4\xc2\xa2\x4e\x17\x76\x65\x30\x2f\x2b\xbc\x9b\xe2\x13\x4b\x9b\x7d\xd8\x0a\x99\x69\x8c\xe7\x47\xc4\x14\x21\xf6\xc7\x14\x32\xf0\x11\x84\x54\xa0\x29\x8f\x55\x8c\xb4\x48\x2b\x33\x6f\x06\x4e\x44\x53\x45\x22\xe6\x52\x9c\xe8\x3a\x9b\x4d\x9e\xf5\xcb\xba\x5e\x77\x06\x9e\x13\x2b\xf6\x18\x90\xb3\xb8\x5d\x2f\xcd\xce\x56\xb5\xd7\xf5\x00\x20\xcf\xe0\x46\x6b\x3e\xc4\x99\xd9\x47\x08\x1d\x04\x98\x44\x3c\x64\xd1\x0c\x4f\x34\xbb\x69\x0d\x36\xb5\xad\x0a\xb3\xb2\x00\x5e\x58\xe4\x8f\x59\xd4\xbd\xbb\x7f\xe8\x9d\x33\xf8\x85\xd6\x65\x9d\x37\x79\x59\x24\x5b\x0c\xa5\xde\xdf\x6e\x12\x91\xcf\xdb\x07\x5d\x86\x76\xff\x81\xbe\xaf\x80\x76\xe6\xcb\x2d\x46\x81\x64\xea\xc2\xa7\x1c\x92\xf2\x26\x5d\x5c\x0b\xa5\x95\x35\x8d\xcd\x12\xd3\x40\xf1\x90\x62\xc5\xc2\x09\x5e\xb9\x1a\x4b\xad\xda\x13\xbc\x49\x41\x47\xc8\xeb\xfd\xb5\xaa\x05\x7f\xd6\x04\x2e\x86\x34\x45\xbe\x4d\x0e\xbd\x26\xc7\x42\xa4\x38\xed\x1c\x5d\xb7\xec\x0d\xbc\x9f\x00\x00\x00\xff\xff\xce\x06\x3b\x42\x92\x02\x00\x00")

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
