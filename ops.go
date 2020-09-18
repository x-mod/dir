package dir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

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

func (d *Dir) IsEmpty(elems ...string) (bool, error) {
	return afero.IsEmpty(d.fs, d.Path(elems...))
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
