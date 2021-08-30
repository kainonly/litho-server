package controller

type Index struct {
	*Dependency
}

func NewIndex(d *Dependency) *Index {
	return &Index{
		Dependency: d,
	}
}
