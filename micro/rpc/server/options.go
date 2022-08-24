package server

var defaultServerOptions = serverOptions{}

type serverOptions struct {
	openssl  bool
	certFile string
	keyFile  string

	wraps []HandlerWrapper
}

type ServerOption interface {
	apply(*serverOptions)
}

type EmptyServerOption struct{}

func (EmptyServerOption) apply(*serverOptions) {}

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

func SetRSAKey(cartfile, keyfile string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.certFile = cartfile
		o.keyFile = keyfile

		if o.certFile != "" && o.keyFile != "" {
			o.openssl = true
		}
	})
}

func WithHandlerWrap(hw ...HandlerWrapper) ServerOption {
	return newFuncServerOption(func(options *serverOptions) {
		options.wraps = append(options.wraps, hw...)
	})
}
