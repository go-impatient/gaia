package request

// Method ...
func Method(method string) *Request {
	return New().Method(method)
}

// Url ...
func Url(url string) *Request {
	return New().Url(url)
}

// Get ...
func Get(url string) *Request {
	return New().Get(url)
}

// Post ...
func Post(url string) *Request {
	return New().Post(url)
}

// Put ...
func Put(url string) *Request {
	return New().Put(url)
}

// Delete ...
func Delete(url string) *Request {
	return New().Delete(url)
}
