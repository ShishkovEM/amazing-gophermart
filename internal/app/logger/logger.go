package logger

import (
	"io"
	"net/http"

	"github.com/izumin5210/httplogger"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

func New(out io.Writer) plugin.Plugin {
	return create(func(parent http.RoundTripper) http.RoundTripper {
		return httplogger.NewRoundTripper(out, parent)
	})
}

func create(transportFn func(parent http.RoundTripper) http.RoundTripper) plugin.Plugin {
	return plugin.NewRequestPlugin(func(c *context.Context, h context.Handler) {
		c.Client.Transport = transportFn(c.Client.Transport)
		h.Next(c)
	})
}
