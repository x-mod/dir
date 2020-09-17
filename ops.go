package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

type Dir struct {
	root       string
	folders    []string
	createMode os.FileMode
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

func Folder(folder string) Option {
	return func(d *Dir) {
		d.folders = append(d.folders, folder)
	}
}

func New(opts ...Option) *Dir {
	d := &Dir{
		root:       ".",
		folders:    []string{},
		createMode: 0700,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dir) Open() error {
	if _, err := os.Stat(d.root); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(d.root, d.createMode); err != nil {
				return fmt.Errorf("root mkdir: %w", err)
			}
		}
	}
	abs, err := filepath.Abs(d.root)
	if err != nil {
		return fmt.Errorf("root abs: %w", err)
	}
	d.root = abs
	for _, folder := range d.folders {
		if err := os.MkdirAll(filepath.Join(d.root, folder), d.createMode); err != nil {
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
	if err := os.MkdirAll(dest, d.createMode); err != nil {
		return "", err
	}
	return dest, nil
}

func (d *Dir) Path(elems ...string) string {
	ds := []string{d.root}
	ds = append(ds, elems...)
	return filepath.Join(ds...)
}
