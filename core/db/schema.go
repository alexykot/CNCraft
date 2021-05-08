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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\x51\x6f\x9b\x3e\x14\xc5\xdf\xf9\x14\x57\x7d\x0a\x52\xff\x7f\x6d\x7b\x9a\xd4\x27\x17\x6e\x14\x6b\x60\x67\xc6\xde\xd2\xbd\x20\x0a\xde\x6a\xa9\xb1\x23\x63\xb6\x66\x9f\x7e\x0a\x01\x16\x25\x65\xeb\xea\x47\xee\xf9\x5d\xfb\x1e\xce\x4d\x04\x12\x89\x50\x24\x2b\xcc\x09\xd0\x25\x30\x2e\x01\x37\xb4\x90\x05\xd4\xb6\xf6\xd5\xd7\x70\x13\x0d\x22\xdc\x48\x64\x05\xe5\xec\x4c\x77\xd5\x75\xa6\xf9\xcf\xb5\xed\xee\xea\x26\x1a\xc5\x92\xdc\x66\x38\xb6\xf8\x7f\xf7\x58\xed\xb5\x6f\xa3\x45\x04\x00\x60\x1a\x38\x3d\x4a\xd1\x14\x66\xce\xe1\x1a\xa6\xb2\xec\xba\x07\x6b\x67\x6d\xf9\x9b\xfe\x23\x38\x41\x5d\xab\xbd\xad\xb6\x7a\x28\x7c\x22\x22\x59\x11\xb1\x78\xfb\xee\x7d\x3c\x7f\x5b\x4f\xee\x5c\x6b\x82\x71\xb6\x7c\xea\xab\x29\x57\x87\x99\xd6\x02\x13\xda\xdb\x30\xfb\xce\x09\xdc\xbf\x16\xfc\xf9\x8f\xe0\xbe\xfa\x71\x5a\x5c\x66\x9c\xc8\x17\x59\xba\x33\xa1\x7e\x78\x0d\xe8\x6c\xf9\xcd\xbb\xce\x0e\x7f\xe3\x96\xf3\x6c\x86\x9b\x40\x48\x71\x49\x54\x26\x41\x0a\x85\x83\xc7\x75\xe7\xbd\xb6\xa1\x7c\x70\xe1\xbe\xf2\x50\xe4\x24\xcb\x28\x7b\xf6\x09\x17\x6d\xde\x8c\x3d\xbc\xae\x82\x6e\xca\x2a\xf4\x3a\x49\x73\x2c\x24\xc9\xd7\xf0\x99\xca\x15\x57\xb2\xff\x02\x5f\x38\xc3\xb3\x19\xd6\x82\xe6\x44\xdc\xc1\x07\xbc\x83\x85\x69\xe2\x28\x9e\xc2\xae\x18\xfd\xa8\x10\x28\x4b\x71\x03\xe6\xa9\x1c\x12\x5c\x4e\x69\xe2\xec\x3c\xdd\xb0\x18\x8b\xf1\xdc\x1e\x18\xfb\x5d\xdb\xe0\xfc\x7e\xd8\x84\x23\x79\x8c\x74\x9f\x66\x81\x4b\x14\xc8\x12\x2c\x2e\xbb\x9b\x26\x3e\xdc\x9a\x62\x86\x12\x21\x21\x45\x42\x52\x3c\x0e\xd2\x3e\xba\x50\xda\x6e\x7b\xaf\x4f\x2c\x7c\xce\xaf\x83\xd8\x04\xbd\x1d\xb7\xe8\x65\xe2\xda\x75\x36\xfc\x45\x7c\x69\xe8\x34\xdc\xf5\xe9\xfb\x7a\x93\x7f\x05\x00\x00\xff\xff\x94\xc9\x09\x64\x77\x04\x00\x00")

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
