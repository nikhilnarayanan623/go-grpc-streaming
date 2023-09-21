package file

import (
	"io/fs"
	"os"
)

/**
*  Abstract os file system
**/

// handler to abstract all file system functionalities needed
type Handler interface {
	Create(name string) (File, error)
	MkdirAll(path string, perm fs.FileMode) error
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) MkdirAll(path string, perm fs.FileMode) error {

	return os.MkdirAll(path, perm)
}

func (h *handler) Create(name string) (File, error) {

	// create a fie of os
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	// return it as our abstracted file
	return &file{
		fl: f,
	}, nil
}

// abstraction of each file
type File interface {
	Write(data []byte) (int, error)
	Close() error
}

type file struct {
	fl *os.File
}

func (f *file) Write(data []byte) (int, error) {
	// write the data on actual file
	return f.fl.Write(data)
}

func (f *file) Close() error {
	// close the actual file
	return f.fl.Close()
}
