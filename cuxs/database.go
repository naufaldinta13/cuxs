// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cuxs

import (
	"fmt"
	"time"

	"github.com/naufaldinta13/cuxs/common/log"
	"github.com/naufaldinta13/cuxs/orm"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DbSetup registering database connection
// by reading database config variables.
func DbSetup() error {
	orm.Debug = IsDebug()
	orm.DebugLog = log.Log
	orm.DefaultTimeLoc = time.Local
	orm.DefaultRelsDepth = 3

	ds := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", Config.DbUser, Config.DbPassword, Config.DbHost, Config.DbName, "charset=utf8&loc=Asia%2FJakarta")
	return orm.RegisterDataBase("default", Config.DbEngine, ds)
}
