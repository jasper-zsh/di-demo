package demo

import "di-demo/di"

type StructA struct {
	B *StructB
}

func (s *StructA) Requires() []di.Dependency {
	return []di.Dependency{}
}

func (s *StructA) Inject(c *di.Injector) error {
	b, err := c.Get(&StructB{})
	if err != nil {
		return err
	}
	s.B = b.(*StructB)
	return nil
}

type StructB struct {
	A *StructA
}

func (s *StructB) Requires() []di.Dependency {
	return []di.Dependency{&StructA{}}
}

func (s *StructB) Inject(container *di.Injector) error {
	a, err := container.Get(&StructA{})
	if err != nil {
		return err
	}
	s.A = a.(*StructA)
	return nil
}

type StructC struct {
	A   *StructA
	B   *StructB
	Foo string
}

func (s *StructC) Requires() []di.Dependency {
	return []di.Dependency{&StructA{}, &StructB{}}
}

func (s *StructC) Inject(container *di.Injector) error {
	a, err := container.Get(&StructA{})
	if err != nil {
		return err
	}
	b, err := container.Get(&StructB{})
	if err != nil {
		return err
	}
	s.A = a.(*StructA)
	s.B = b.(*StructB)

	return nil
}
