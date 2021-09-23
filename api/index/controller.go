package index

func (x *Controller) Index() interface{} {
	return x.Service.Version()
}
