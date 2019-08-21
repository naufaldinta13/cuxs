// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cuxs

import (
	"net/http"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Context is custom echo.Context
// has defined as middleware.
type Context struct {
	echo.Context
	responseFormat *ResponseFormat
	ResponseData   interface{}
}

// NewContext new instances of context
func NewContext(c echo.Context) *Context {
	return &Context{c, NewResponse(), nil}
}

// Decrypt helper for decripting param it will return error not found
// if failed to decrypt or the value is 0
func (c *Context) Decrypt(param string) (id int64, e error) {
	id, e = common.Decrypt(c.Param(param))

	if e != nil || id == 0 {
		e = echo.ErrNotFound
	}

	return
}

// Data set data and total into response format
func (c *Context) Data(data interface{}, total ...int64) {
	c.responseFormat.SetData(data, total...)
}

// Files set data and total into response format
func (c *Context) Files(file string) {
	c.responseFormat.SetFile(file)
}

// Failure set response format errors
// its equal with validation errors.
func (c *Context) Failure(fail ...string) {
	c.responseFormat.Errors = map[string]string{fail[0]: fail[1]}
}

// Serve response json data with data that already collected
// if error is not nill will returning error responses.
func (c *Context) Serve(e error) (err error) {
	c.responseFormat.Code = http.StatusOK
	if e != nil {
		c.responseFormat.SetError(e)
	} else {
		if c.ResponseData != nil {
			c.responseFormat.SetData(c.ResponseData)
		}
	}

	if len(c.responseFormat.Errors) > 0 {
		c.responseFormat.Status = HTTPResponseFail
		c.responseFormat.Code = http.StatusUnprocessableEntity
		c.responseFormat.Message = http.StatusText(http.StatusUnprocessableEntity)
		c.responseFormat.Data = nil
		c.responseFormat.Total = 0
	}

	if c.Request().Method == echo.HEAD || c.Request().Method == echo.OPTIONS {
		err = c.NoContent(http.StatusNoContent)
	} else {
		err = c.JSON(c.responseFormat.Code, c.responseFormat)
	}

	c.responseFormat.reset()

	return
}

// RequestQuery set query param into orm so the repository
// can use the data.
func (c *Context) RequestQuery() *orm.RequestQuery {
	rq := &orm.RequestQuery{
		Offset:     0,
		Limit:      -1,
		Conditions: make([]map[string]string, 0),
	}

	rq = rq.ReadFromContext(c.QueryParams())
	if c.QueryParam("export") == "1" {
		rq.Limit = -1
	}

	return rq
}

// JwtUsers get a user sessions that having jwt token in
// request header and checked again the model.
func (c *Context) JwtUsers(model jwtUser) interface{} {
	if u := c.Get("user"); u != nil {
		s := u.(*jwt.Token)
		c := s.Claims.(jwt.MapClaims)
		id := int64(c["id"].(float64))

		if users, err := model.GetUser(id); err == nil {
			return users
		}
	}

	return nil
}

// jwtUser model user jwt token interface
// to check is the id given valid as users.
type jwtUser interface {
	GetUser(int64) (interface{}, error)
}
