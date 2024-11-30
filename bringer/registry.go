package bringer

type Factory = func(opts ...Option) Bringer

var bringers = map[string](Factory){}

func Register(schema string, f Factory) {
	bringers[schema] = f
}
