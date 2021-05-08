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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\xc1\x6f\xd3\x30\x14\xc6\xef\xfd\x2b\x9e\x76\x6a\xa4\x81\x80\x13\xd2\x4e\x5e\xf2\xaa\x5a\x24\x76\x71\x6c\xe8\xb8\x44\x59\x62\x98\xa5\xd5\xae\x1c\x07\x56\xfe\x7a\x94\x34\x09\x55\xb7\xc0\x36\x1f\xf3\xbe\xdf\x7b\x7e\x5f\x3e\xc7\x02\x89\x44\xc8\xe3\x35\x66\x04\xe8\x0a\x18\x97\x80\x5b\x9a\xcb\x1c\x2a\x5b\xf9\xf2\x7b\xb8\x5a\x0c\x22\xdc\x4a\x64\x39\xe5\xec\x4c\x77\xd1\xb6\xa6\x7e\xe3\x9a\x66\x7f\x71\xb5\x18\xc5\x92\x5c\xa7\x38\xb6\x78\xbb\xbf\x2f\x0f\xda\x37\x8b\xe5\x02\x00\xc0\xd4\x70\x7a\x94\xa2\x09\xcc\x9c\x6e\x0c\x53\x69\x7a\xd9\x83\x95\xb3\xb6\xf8\x4b\xff\x13\x9c\xa0\xb6\xd1\xde\x96\x3b\x3d\x14\xbe\x10\x11\xaf\x89\x58\xbe\xff\xf0\x31\x9a\x9f\xd6\x93\x7b\xd7\x98\x60\x9c\x2d\x1e\xfa\x6a\xc2\x55\xb7\xd3\x46\x60\x4c\x7b\x1b\x66\xef\x39\x81\x87\xd7\x82\xbf\x5f\x08\x1e\xca\x5f\xa7\xc5\x55\xca\x89\x7c\x96\xa5\x7b\x13\xaa\xbb\xd7\x80\xce\x16\x3f\xbc\x6b\xed\xf0\x37\xae\x39\x4f\x67\xb8\x09\x84\x04\x57\x44\xa5\x12\xa4\x50\x38\x78\x5c\xb5\xde\x6b\x1b\x8a\x3b\x17\x6e\x4b\x0f\x94\xcd\x4d\x7f\xa2\xcd\xbb\xb1\x87\xd7\x65\xd0\x75\x51\x86\x5e\x27\x69\x86\xb9\x24\xd9\x06\xbe\x52\xb9\xe6\x4a\xf6\x5f\xe0\x1b\x67\x78\xb6\xc3\x46\xd0\x8c\x88\x1b\xf8\x84\x37\xb0\x34\x75\xb4\x88\xa6\xb0\x2b\x46\x3f\x2b\x04\xca\x12\xdc\x82\x79\x28\x86\x04\x17\x53\x9a\x38\x3b\x4f\x37\x2c\xc7\x62\x34\xf7\x0e\x8c\xfd\xa9\x6d\x70\xfe\x30\xbc\x84\x23\x79\x8c\x74\x9f\x66\x81\x2b\x14\xc8\x62\xcc\x1f\x77\x37\x75\xd4\x4d\x4d\x30\x45\x89\x10\x93\x3c\x26\x09\x1e\x17\x69\xee\x5d\x28\x6c\xbb\xbb\xd5\x1e\xf2\x8c\xa4\x69\xe7\xe3\x53\x7e\x75\x62\x13\xf4\x6e\x7c\x45\xcf\x13\x57\xae\xb5\xe1\x3f\xe2\xc7\x86\x4e\xcb\x5d\x9e\xde\xaf\x37\xf9\x4f\x00\x00\x00\xff\xff\x69\x4d\x6a\x43\x77\x04\x00\x00")

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
