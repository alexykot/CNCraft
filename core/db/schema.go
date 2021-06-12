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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\xc1\x6f\x9b\x30\x14\xc6\xef\xfc\x15\x4f\x3d\x25\x52\x36\x6d\x3b\x4d\xea\xc9\x85\x17\xc5\x1a\x98\xcc\xd8\x5b\xba\x0b\xa2\xe0\xad\x96\x1a\x3b\x32\x66\x6b\xf6\xd7\x4f\x10\x60\x28\x69\xda\xb4\x3e\xf2\xbe\xdf\x7b\xf6\xc7\xfb\x42\x8e\x44\x20\x64\xe1\x0a\x13\x02\x74\x09\x2c\x15\x80\x1b\x9a\x89\x0c\x4a\x53\xba\xe2\xa7\xbf\x0e\x7a\x11\x6e\x04\xb2\x8c\xa6\xec\x48\x77\xd5\x34\xba\x7a\x67\xeb\x7a\x77\x75\x1d\x0c\x62\x41\x6e\x62\x1c\x5a\xbc\xdf\x3d\x14\x7b\xe5\xea\x60\x16\x00\x00\xe8\x0a\xa6\x47\x4a\x1a\xc1\x99\xd3\x8e\x61\x32\x8e\x17\x1d\x58\x5a\x63\xf2\xff\xf4\xb3\xe0\x08\x55\x7a\xab\x4c\xad\x6d\x4f\x5e\x3e\xad\xa9\x95\x33\xc5\x56\xf5\xc5\x6f\x84\x87\x2b\xc2\x67\x1f\x3f\x7d\x9e\x3f\x0f\xee\x6c\xad\x7d\x3b\xf0\xb1\x2b\x46\xa9\x6c\xbd\x58\x73\x0c\x69\x67\xdf\xcb\xe0\xfe\xad\xe0\xdf\x57\x82\xfb\xe2\xcf\xb4\xb8\x8c\x53\x22\x2e\x32\x67\xa7\x7d\x79\xff\x16\xd0\x9a\xfc\x97\xb3\x8d\xe9\xff\xe2\x4d\x9a\xc6\x67\xb8\x11\x84\x08\x97\x44\xc6\x02\x04\x97\xb8\x08\x0e\xab\xd0\x38\xa7\x8c\xcf\xef\xad\xbf\x2b\x1c\x64\x09\x89\x63\xca\x9e\xbc\xc2\x49\x9b\x0f\x43\x0f\xa7\x0a\xaf\xaa\xbc\xf0\x9d\x4e\xd0\x04\x33\x41\x92\x35\x7c\xa7\x62\x95\x4a\xd1\x7d\x81\x1f\x29\xc3\xa3\x37\xac\x39\x4d\x08\xbf\x85\x2f\x78\x0b\x33\x5d\xcd\x83\xf9\x18\x12\xc9\xe8\x57\x89\x40\x59\x84\x1b\xd0\x8f\x79\xbf\xf9\xf9\xb8\x4c\x29\x3b\x4e\x05\xcc\x86\xe2\xfc\x5c\x7e\xb4\xf9\xad\x8c\xb7\x6e\xdf\x27\xe8\x40\x4e\x16\x9a\xe3\x12\x39\xb2\x10\xb3\xd3\xee\xba\x9a\xb7\x53\x23\x8c\x51\x20\x84\x24\x0b\x49\x84\x87\x87\xd4\x0f\xd6\xe7\xa6\xd9\xde\xa9\x89\x85\x4f\xf9\xd5\x8a\xb5\x57\xdb\x21\x7d\x97\x89\x4b\xdb\x18\xff\x82\xf8\xd4\xd0\xf1\x71\x8b\xe9\xfd\x3a\x93\xff\x05\x00\x00\xff\xff\xcd\x2b\x00\x43\xaf\x04\x00\x00")

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
