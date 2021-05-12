package dir

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/spf13/afero"
)

var EmbedFS embed.FS

type Dir struct {
	root       string
	folders    []string
	createMode os.FileMode
	fs         afero.Fs
}

type Option func(*Dir)

func CreateMode(mode os.FileMode) Option {
	return func(d *Dir) {
		d.createMode = mode
	}
}

func Root(root string) Option {
	return func(d *Dir) {
		d.root = root
	}
}

func Folder(elems ...string) Option {
	return func(d *Dir) {
		d.folders = append(d.folders, filepath.Join(elems...))
	}
}

func New(opts ...Option) *Dir {
	d := &Dir{
		root:       ".",
		folders:    []string{},
		createMode: 0700,
		fs:         afero.NewOsFs(),
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dir) Open() error {
	abs, err := filepath.Abs(d.root)
	if err != nil {
		return fmt.Errorf("root abs: %w", err)
	}
	d.root = abs

	exist, err := afero.DirExists(d.fs, d.root)
	if err != nil {
		return err
	}
	if !exist {
		if err := d.fs.MkdirAll(d.root, d.createMode); err != nil {
			return fmt.Errorf("root mkdir: %w", err)
		}
	}
	for _, folder := range d.folders {
		if err := d.fs.MkdirAll(filepath.Join(d.root, folder), d.createMode); err != nil {
			return fmt.Errorf("folder mkdir: %w", err)
		}
	}
	return nil
}

func (d *Dir) String() string {
	return d.root
}

func (d *Dir) Mkdir(elems ...string) (string, error) {
	ds := []string{d.root}
	ds = append(ds, elems...)
	dest := filepath.Join(ds...)
	if err := d.fs.MkdirAll(dest, d.createMode); err != nil {
		return "", err
	}
	return dest, nil
}

func (d *Dir) Path(elems ...string) string {
	ds := []string{d.root}
	ds = append(ds, elems...)
	return filepath.Join(ds...)
}

func (d *Dir) DirExists(elems ...string) (bool, error) {
	return afero.DirExists(d.fs, d.Path(elems...))
}

func (d *Dir) Exists(elems ...string) (bool, error) {
	return afero.Exists(d.fs, d.Path(elems...))
}

func (d *Dir) IsEmpty(elems ...string) (bool, error) {
	return afero.IsEmpty(d.fs, d.Path(elems...))
}

func (d *Dir) Remove(elems ...string) error {
	return d.fs.Remove(d.Path(elems...))
}

func (d *Dir) RemoveAll(elems ...string) error {
	return d.fs.RemoveAll(d.Path(elems...))
}

func (d *Dir) Folders(elems ...string) ([]string, error) {
	infos, err := afero.ReadDir(d.fs, d.Path(elems...))
	if err != nil {
		return nil, err
	}
	dirs := []string{}
	for _, info := range infos {
		if info.IsDir() {
			dirs = append(dirs, info.Name())
		}
	}
	return dirs, nil
}

func (d *Dir) Files(elems ...string) ([]string, error) {
	infos, err := afero.ReadDir(d.fs, d.Path(elems...))
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
	}
	return files, nil
}

func (d *Dir) SymLink(source string, target string) error {
	if symlink, ok := d.fs.(afero.Symlinker); ok {
		return symlink.SymlinkIfPossible(source, target)
	}
	return fmt.Errorf("unsupport symlinker")
}

func (d *Dir) Unlink(elems ...string) error {
	return syscall.Unlink(d.Path(elems...))
}

func (d *Dir) Stat(elems ...string) (os.FileInfo, error) {
	return d.fs.Stat(d.Path(elems...))
}

func (d *Dir) Embed(elems ...string) error {
	if _, err := d.Mkdir(elems...); err != nil {
		return err
	}
	entries, err := EmbedFS.ReadDir(path.Join(elems...))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		dst := []string{}
		dst = append(dst, elems...)
		dst = append(dst, entry.Name())
		if entry.IsDir() {
			if err := d.Embed(dst...); err != nil {
				return err
			}
			continue
		}
		exist, err := d.Exists(dst...)
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		wr, err := os.Create(d.Path(dst...))
		if err != nil {
			return fmt.Errorf("file <%s> create: %w", path.Join(dst...), err)
		}
		rd, err := EmbedFS.Open(path.Join(dst...))
		if err != nil {
			return fmt.Errorf("file <%s> open: %w", path.Join(dst...), err)
		}
		if _, err := io.Copy(wr, rd); err != nil {
			return fmt.Errorf("file <%s> write: %w", path.Join(dst...), err)
		}
	}
	return nil
}
