// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package tidb

import (
	"fmt"

	"github.com/ngaut/log"
	mysql "github.com/pingcap/tidb/mysqldef"
	"github.com/pingcap/tidb/util/errors"
	"github.com/pingcap/tidb/util/errors2"
)

// Bootstrap initiates system DB for a store.
func bootstrap(s *session) {
	// Create a test database.
	_, err := s.Execute("CREATE DATABASE IF NOT EXISTS test")
	if err != nil {
		log.Fatal(err)
	}

	//  Check if system db exists.
	_, err = s.Execute(fmt.Sprintf("USE %s;", mysql.SystemDB))
	if err == nil {
		// We have already finished bootstrap.
		return
	} else if !errors2.ErrorEqual(err, errors.ErrDatabaseNotExist) {
		log.Fatal(err)
	}
	_, err = s.Execute(fmt.Sprintf("CREATE DATABASE %s;", mysql.SystemDB))
	if err != nil {
		log.Fatal(err)
	}
	initUserTable(s)
}

func initUserTable(s *session) {
	_, err := s.Execute(mysql.CreateUserTable)
	if err != nil {
		log.Fatal(err)
	}
	// Insert a default user with empty password.
	_, err = s.Execute(`INSERT INTO mysql.user VALUES ("localhost", "root", ""), ("127.0.0.1", "root", "");`)
	if err != nil {
		log.Fatal(err)
	}
}
