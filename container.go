package dic

import (
	"errors"
)

var (
	ErrDependencyAlreadyExists = errors.New("Dependency already exists.")
	ErrDependencyNotFound      = errors.New("Dependency not found.")
)

type entry struct {
	loader func(c *container) interface{}
	cache  interface{}
	loaded bool
}

type container struct {
	entries map[string]*entry
}

func New() *container {
	return &container{
		entries: make(map[string]*entry),
	}
}

func (c *container) Set(key string, loader func(c *container) interface{}) error {
	_, ok := c.entries[key]
	if ok {
		return ErrDependencyAlreadyExists
	}

	c.entries[key] = &entry{
		loader: loader,
		cache:  nil,
		loaded: false,
	}

	return nil
}

func (c *container) Get(key string) (interface{}, error) {
	dep, ok := c.entries[key]
	if !ok {
		return nil, ErrDependencyNotFound
	}

	if !dep.loaded {
		dep.cache = dep.loader(c)
		dep.loaded = true
	}

	return dep.cache, nil
}
