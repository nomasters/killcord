package shell

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/errgo.v1"
)

type shellVFS struct {
	root  string
	shell *Shell
}

func (svfs shellVFS) Open(name string) (http.File, error) {
	p := filepath.Join(svfs.root, name)
	obj, err := svfs.shell.FileList(p)
	if err != nil {
		return nil, errgo.Notef(err, "vfs: failed to list path: %s", p)
	}
	switch obj.Type {
	case "Directory":
		return &shellvfs_dir{
			&shellvfs_dirInfo{p, obj.Links},
			0,
		}, nil

	case "File":
		finfo := &shellvfs_fileInfo{p,
			&UnixLsLink{Hash: obj.Hash, Name: name, Size: obj.Size, Type: obj.Type},
		}
		body, err := svfs.shell.Cat(obj.Hash)
		if err != nil {
			return nil, errgo.Notef(err, "vfs: failed to list path: %s", p)
		}
		return &shellvfs_file{finfo, body}, nil
	default:
		return nil, errgo.Newf("vfs: unhandled object type on <%s>: %s", p, obj.Type)
	}
}

func NewVFS(root string, shell *Shell) http.FileSystem {
	if shell == nil {
		shell = NewShell("localhost:5001")
	}
	return &shellVFS{root, shell}
}

type shellvfs_fileInfo struct {
	path string
	obj  *UnixLsLink
}

func (f *shellvfs_fileInfo) Readdir(int) ([]os.FileInfo, error) {
	return nil, errgo.Newf("cannot Readdir from file %s", f.path)
}

func (f *shellvfs_fileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *shellvfs_fileInfo) Name() string       { return f.obj.Name }
func (f *shellvfs_fileInfo) Size() int64        { return int64(f.obj.Size) }
func (f *shellvfs_fileInfo) Mode() os.FileMode  { return 0444 }
func (f *shellvfs_fileInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (f *shellvfs_fileInfo) IsDir() bool        { return false }
func (f *shellvfs_fileInfo) Sys() interface{}   { return nil }

type shellvfs_file struct {
	*shellvfs_fileInfo
	io.ReadCloser
}

func (d shellvfs_file) Seek(offset int64, whence int) (int64, error) {
	return 0, errgo.Newf("seek not supported - please open an issue")
}

type shellvfs_dirInfo struct {
	name    string
	entries []*UnixLsLink
}

func (d shellvfs_dirInfo) Read(p []byte) (n int, err error) {
	return 0, errgo.Newf("vfs: <%s> is a Directory", d.name)
}

func (d shellvfs_dirInfo) Stat() (os.FileInfo, error) { return d, nil }
func (d shellvfs_dirInfo) Close() error               { return nil }

func (d shellvfs_dirInfo) IsDir() bool        { return true }
func (d shellvfs_dirInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (d shellvfs_dirInfo) Mode() os.FileMode  { return 0 }
func (d shellvfs_dirInfo) Name() string       { return d.name }
func (d shellvfs_dirInfo) Size() int64        { return int64(len(d.entries)) }
func (d shellvfs_dirInfo) Sys() interface{}   { return nil }

type shellvfs_dir struct {
	*shellvfs_dirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *shellvfs_dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == os.SEEK_SET {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.shellvfs_dirInfo.name)
}

func (d *shellvfs_dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.shellvfs_dirInfo.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.shellvfs_dirInfo.entries)-d.pos {
		count = len(d.shellvfs_dirInfo.entries) - d.pos
	}
	e := d.shellvfs_dirInfo.entries[d.pos : d.pos+count]
	d.pos += count

	return wrapFileInfo(e), nil
}

func wrapFileInfo(in []*UnixLsLink) []os.FileInfo {
	out := make([]os.FileInfo, len(in))
	for i, link := range in {
		out[i] = &shellvfs_fileInfo{"", link}
	}
	return out
}
