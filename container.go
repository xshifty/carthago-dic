package dic

import (
	"errors"
	"sync"
)

var ErrDependencyAlreadyExists = errors.New("Dependency already exists.")
var ErrDependencyNotFound = errors.New("Dependency not found.")

type Container struct {
	mutex   *sync.Mutex
	entries map[string]*entry
}

type entry struct {
	loader func(c *Container) interface{}
	cache  interface{}
	loaded bool
}

func NewContainer() *Container {
	return &Container{&sync.Mutex{}, make(map[string]*entry)}
}

func (c *Container) Set(key string, loader func(c *Container) interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	c.mutex.Lock()
	defer c.mutex.Unlock()

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

func (c *Container) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.entries, key)
}
