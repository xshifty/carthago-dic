package dic

import (
	"errors"
)

var ErrDependencyAlreadyExists = errors.New("Dependency already exists.")
var ErrDependencyNotFound = errors.New("Dependency not found.")

type Container struct {
	entries map[string]*entry
}

type entry struct {
	loader func(c *Container) interface{}
	cache  interface{}
	loaded bool
}

func New() *Container {
	return &Container{
		entries: make(map[string]*entry),
	}
}

func (c *Container) Set(key string, loader func(c *Container) interface{}) error {
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

func (c *Container) Get(key string) (interface{}, error) {
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
