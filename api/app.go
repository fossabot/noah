/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 */

package api

import (
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/kataras/iris"
	"net/http"
)

func StartApp(sctx *system.SContext) {
	app := iris.Default()

	app.Get("/nodes", func(ctx iris.Context) {
		if nodes, err := sctx.GetNodes(); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct{
				Error string
			}{
				Error: err.Error(),
			})
		} else {
			ctx.JSON(nodes)
		}
	})

	app.Put("/nodes", func(ctx iris.Context) {
		node := system.NNode{}
		if err := ctx.ReadJSON(&node); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct{
				Error string
			}{
				Error: err.Error(),
			})
		} else if err := sctx.AddNode(node); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct {
				Error string
			}{
				Error: err.Error(),
			})
		} else {
			ctx.StatusCode(200)
			ctx.JSON(iris.Map{
				"Message": "Node has been created.",
			})
		}
	})

	app.Get("/settings", func(ctx iris.Context) {
		if settings, err := sctx.GetSettings(); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct{
				Error string
			}{
				Error: err.Error(),
			})
		} else {
			ctx.JSON(settings)
		}
	})

	app.Build()
	srv := &http.Server{Handler: app, Addr: ":8080"} // you have to set Handler:app and Addr, see "iris-way" which does this automatically.
	// http://localhost:8080/
	// http://localhost:8080/mypath
	println("Start a server listening on http://localhost:8080")
	srv.ListenAndServe() // same as app.Run(iris.Addr(":8080"))


	// listen and serve on http://0.0.0.0:8080.
	//app.Run(iris.Addr(fmt.Sprintf(":%d", sctx.Flags.HTTPPort)))
}