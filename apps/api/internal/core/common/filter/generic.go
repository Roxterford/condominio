package filter

type generic struct {
	spec Spec
}

func (g generic) FilterSpec() Spec {
	return g.spec
}

func G(s Spec) Filterable {
	return generic{s}
}
