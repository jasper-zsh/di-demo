package di

type Dependency interface {
	Requires() []Dependency
	Inject(injector *Injector) error
}
