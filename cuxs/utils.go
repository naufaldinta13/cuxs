// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cuxs

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"time"

	"git.qasico.com/cuxs/common/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Authorized returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For empty token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func Authorized() echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = JwtKey()
	return middleware.JWTWithConfig(c)
}

// JwtKey byte of jwt secret keys.
func JwtKey() []byte {
	return []byte(Config.JwtSecret)
}

// JwtToken make an JWT token keys and values
// the return will become a valid token with
// a life time 72 hours from the time generated.
func JwtToken(k string, v interface{}) (token string) {
	// new instances jwt
	jwts := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := jwts.Claims.(jwt.MapClaims)
	claims[k] = v
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	var e error
	if token, e = jwts.SignedString(JwtKey()); e != nil {
		panic(e)
	}

	return token
}

// HTTPCallTest fungsi http call untuk keperluan testing
func HTTPCallTest(method string, data io.Reader, token string) (ctx *Context, res *httptest.ResponseRecorder) {
	token = "Bearer " + token
	req, _ := http.NewRequest(method, "/", data)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	res = httptest.NewRecorder()
	c := New().NewContext(req, res)
	ctx = NewContext(c)

	return
}

// ListRoutes print all route available, only show on debug mode.
func listRoutes(e *echo.Echo) {
	log.Debug("%0120v", "")
	log.Debug("%-10s | %-50s | %-54s", "METHOD", "URL PATH", "REQ. HANDLER")
	log.Debug("%0120v", "")

	routes := e.Routes()
	sort.Sort(sortByPath(routes))
	for _, v := range routes {
		if v.Path[len(v.Path)-1:] != "*" {
			log.Debug("%-10s | %-50s | %-54s", v.Method, v.Path, filepath.Base(v.Name))
		}
	}
	log.Debug("%0120v", "")
}

// sortByPath Sorting echo.Routes by path
// so it make more pretty when printed on console.
type sortByPath []*echo.Route

func (a sortByPath) Len() int {
	return len(a)
}

func (a sortByPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortByPath) Less(i, j int) bool {
	if a[i].Path < a[j].Path {
		return true
	}
	if a[i].Path > a[j].Path {
		return false
	}
	return a[i].Path < a[j].Path
}
