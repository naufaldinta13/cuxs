// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cuxs

import (
	"github.com/naufaldinta13/cuxs/cuxs/mw"

	"github.com/naufaldinta13/cuxs/common/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New creates an instance of Echo.
func New() *echo.Echo {
	// Make new instaces echo
	e := echo.New()

	// custom middleware
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{c, NewResponse(), nil}
			return h(cc)
		}
	})

	// using custom binder
	e.Binder = &binder{}
	e.HTTPErrorHandler = HTTPErrorHandler

	// registering middleware
	e.Use(mw.HTTPLogger())
	e.Use(middleware.Gzip(), middleware.CORS())
	e.Use(middleware.Secure(), middleware.Recover())

	return e
}

// StartServer starting echo servers
func StartServer(e *echo.Echo) {
	if IsDebug() {
		listRoutes(e)
	}

	// starting webserver
	if err := e.Start(Config.Host); err != nil {
		log.Error(err)
	}
}

// HTTPErrorHandler invokes the default HTTP error handler.
func HTTPErrorHandler(err error, c echo.Context) {
	if !c.Response().Committed {
		ctx, ok := c.(*Context)
		if !ok {
			ctx = NewContext(c)
		}

		ctx.Serve(err)
	}
}
