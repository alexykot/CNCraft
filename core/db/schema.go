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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x94\x5f\x6f\xb2\x30\x14\xc6\xef\xf9\x14\x27\x5e\x69\xe2\xfb\x46\x7d\xff\x6c\x89\x57\x15\x8e\xb1\x19\xb4\xae\xb4\x9b\xee\x86\x10\x64\x1b\x89\x03\x53\x4a\xa6\xfb\xf4\x8b\x08\x9a\x19\x8d\x51\x7a\x49\x9f\xe7\x77\x0e\xe7\x34\x8f\x2d\x90\x48\x04\xdf\x9e\xa0\x47\x80\x8e\x81\x71\x09\x38\xa3\xbe\xf4\x21\x4a\x23\x1d\xbe\x9a\xa1\x55\x89\x70\x26\x91\xf9\x94\xb3\x23\x5d\xab\x28\x92\xc5\xaf\x2c\xcf\x57\xad\xa1\x55\x8b\x25\x19\xb9\x58\x23\x7e\xaf\x96\xe1\x26\xd6\xb9\xd5\xb6\x00\x00\x92\x05\x1c\x8e\x52\xd4\x81\x33\x67\x5b\x84\x29\xd7\x85\xa9\xa0\x1e\x11\x73\x78\xc0\x79\xb7\x44\x14\x79\xac\xd3\xf0\x23\x2e\x65\x4f\x44\xd8\x13\x22\xda\xfd\xc1\x7d\xe7\x1c\xa2\x6b\x95\xbe\x55\x96\x27\x26\xc9\xd2\x60\x0d\x00\x0e\x57\xdb\x26\xa7\x02\x6d\x5a\xfe\xd7\x49\xdf\x0f\xdb\xe6\x36\xdb\xd7\x55\xb6\x4d\xf8\x79\xb8\x1a\xbb\x9c\xc8\x0b\xf3\xa9\xaa\x25\x26\x7a\xbf\xde\x96\xa5\xc1\x9b\xce\x8a\xb4\x5c\xca\x88\x73\xf7\xe2\x36\x1c\x1c\x13\xe5\x4a\x90\x42\x61\x35\xd7\xa8\xd0\x3a\x4e\x4d\x90\x2f\x33\x03\x94\x9d\xab\x7c\x02\xd2\xdb\x75\xb1\x35\xf6\x6a\xd1\xad\x80\x7e\x53\xc0\xa0\x29\xe0\x4f\x53\xc0\xdf\xa6\x80\x7f\x4d\x01\xff\x9b\x02\xee\x6e\x05\xec\x5e\x92\x8e\x43\x13\x2f\x82\xd0\x00\x80\xa4\x1e\xfa\x92\x78\x53\x78\xa6\x72\xc2\x95\x2c\xbf\xc0\x0b\x67\xb8\x27\x58\x9d\x43\xea\x28\x46\x1f\x15\x02\x65\x0e\xce\x20\x59\x07\x55\xee\x04\xfb\xb8\xe0\xec\x38\x93\xa0\x5d\x5f\x76\x86\xd6\x77\x00\x00\x00\xff\xff\x26\x71\xde\x54\x0f\x05\x00\x00")

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
