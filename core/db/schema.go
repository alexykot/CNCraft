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

var __001_playersUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x53\xd1\x6e\xd3\x30\x14\x7d\xcf\x57\x5c\xed\xa9\x91\x56\x04\x3c\x21\xed\xc9\x4b\x6e\x55\x0b\xc7\x29\x8e\x03\x2d\x2f\x51\x96\x04\x66\x69\xb5\x2b\xc7\x81\x95\xaf\x47\x31\x6e\x54\xba\x45\x8c\xf9\xd1\xf7\x9c\x7b\xae\xef\x39\x4e\x04\x12\x89\x50\x24\x6b\xcc\x08\xd0\x15\xf0\x5c\x02\x6e\x69\x21\x0b\x68\x74\x63\xeb\x6f\xee\x26\x0a\x20\xdc\x4a\xe4\x05\xcd\xf9\x05\xee\x6a\x18\x54\xbb\x34\x7d\x7f\xb8\xba\x89\x4e\x60\x49\x6e\x19\x9e\x5a\xbc\x39\x3c\xd4\xc7\xce\xf6\xd1\x22\x02\x00\x50\x2d\x9c\x9f\xb2\xa4\x29\xcc\x9c\x51\x86\x97\x8c\x5d\x7b\xe2\xd0\x77\x56\xd7\xfb\x2e\x14\x3f\x13\x91\xac\x89\x58\xbc\x7b\xff\x21\x9e\x27\x7a\xe6\xc1\xf4\xca\x29\xa3\xab\x47\x5f\x4d\xf3\x72\x1c\x6f\x23\x30\xa1\xfe\x45\xb3\x92\x13\xf1\xf8\x5a\xe2\xaf\xff\x24\x1e\xeb\x9f\xe7\xc5\x15\xcb\x89\x7c\xd1\x76\x0e\xca\x35\xf7\xaf\x21\x1a\x5d\x7d\xb7\x66\xd0\xc1\x96\xdb\x3c\x67\x33\xbc\x89\x08\x29\xae\x48\xc9\x24\x48\x51\x62\xd8\x71\x33\x58\xdb\x69\x57\xdd\x1b\x77\x57\x5b\xa0\x7c\x4e\xfd\x99\x36\x6f\x4f\x3d\x6c\x57\xbb\xae\xad\x6a\xe7\x71\x92\x66\x58\x48\x92\x6d\xe0\x0b\x95\xeb\xbc\x94\xfe\x06\xbe\xe6\x1c\x2f\xde\xb0\x11\x34\x23\x62\x07\x1f\x71\x07\x0b\xd5\xc6\x51\x3c\xe5\xb6\xe4\xf4\x53\x89\x40\x79\x8a\x5b\x50\x8f\x55\x08\x63\x35\xa5\x29\xe7\x97\x41\x85\xc5\xa9\x18\xcf\x45\x5a\xe9\x1f\x9d\x76\xc6\x1e\x43\xa8\x97\x4b\xf8\x43\xae\x7c\xbc\x7d\xa8\x05\xae\x50\x20\x4f\xb0\x78\x2a\xa0\xda\x78\x14\x4e\x91\xa1\x44\x48\x48\x91\x90\x14\xc7\x7d\xa0\x10\x5e\x88\x72\x2a\x29\x61\x6c\x17\x2e\x31\x0d\x36\x5f\xaa\xfc\xbd\x88\xfe\xc1\xb8\x4a\x0f\xfb\xbb\xce\x42\x91\x11\xc6\x46\x1f\x9e\xdb\xf7\x08\x56\xae\xdb\x57\xe1\x3b\xbe\x0c\xdc\x98\x41\xbb\x7f\x80\x9f\x1a\x32\xcd\x7c\x7d\x3e\x9f\x37\xe9\x77\x00\x00\x00\xff\xff\xe3\x6f\x71\x08\x82\x04\x00\x00")

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
