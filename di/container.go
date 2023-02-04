package di

import (
	"fmt"
	"reflect"
)

type FactoryFunc func() (Dependency, error)

type DependencyContainer struct {
	factories     map[reflect.Type]FactoryFunc
	injectedCache map[reflect.Type]struct{}
	instanceCache map[reflect.Type]Dependency
	injector      *Injector
}

type Injector struct {
	container *DependencyContainer
}

func (i *Injector) Get(objType Dependency) (Dependency, error) {
	t := reflect.ValueOf(objType).Elem().Type()
	return i.container.instance(t)
}

func NewContainer() *DependencyContainer {
	c := &DependencyContainer{
		factories:     make(map[reflect.Type]FactoryFunc),
		injectedCache: make(map[reflect.Type]struct{}),
		instanceCache: make(map[reflect.Type]Dependency),
	}
	i := &Injector{container: c}
	c.injector = i
	return c
}

func (c *DependencyContainer) instance(objType reflect.Type) (Dependency, error) {
	instance, ok := c.instanceCache[objType]
	if !ok {
		var err error
		factory, ok := c.factories[objType]
		if !ok {
			return nil, fmt.Errorf("factory for %s not registered", objType)
		}
		instance, err = factory()
		if err != nil {
			return nil, err
		}
		c.instanceCache[objType] = instance
		for _, dep := range instance.Requires() {
			_, err = c.instance(reflect.ValueOf(dep).Elem().Type())
			if err != nil {
				return nil, err
			}
		}
	}
	return instance, nil
}

func (c *DependencyContainer) inject(objType reflect.Type) error {
	_, ok := c.injectedCache[objType]
	if !ok {
		instance, ok := c.instanceCache[objType]
		if !ok {
			return fmt.Errorf("instance not found for %s", objType)
		}
		instance.Inject(c.injector)
		c.injectedCache[objType] = struct{}{}
		for _, dep := range instance.Requires() {
			t := reflect.ValueOf(dep).Elem().Type()
			err := c.inject(t)
			if err != nil {
				return fmt.Errorf("failed to inject %s", t)
			}
		}
	}
	return nil
}

func (c *DependencyContainer) Get(objType Dependency) (Dependency, error) {
	val := reflect.ValueOf(objType)
	t := val.Elem().Type()
	instance, err := c.instance(t)
	if err != nil {
		return nil, err
	}
	err = c.inject(t)
	return instance, nil
}

func (c *DependencyContainer) Provide(objType Dependency, factory FactoryFunc) {
	t := reflect.ValueOf(objType).Elem().Type()
	c.factories[t] = factory
}
