package lark

func (x *Controller) RedirectUrl(path string, locale string) []byte {
	u := x.V.Console
	if locale != "" {
		u += "/" + locale
	}
	return []byte(u + path)
}
