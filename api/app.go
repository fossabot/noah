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
 */

package api

import (
	"github.com/kataras/iris"
	"github.com/readystock/noah/db/system"
	"net/http"
)

func StartApp(sctx *system.SContext, addr string) {
	app := iris.Default()

	app.Get("/nodes", func(ctx iris.Context) {
		if nodes, err := sctx.Nodes.GetNodes(); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct {
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
			ctx.JSON(struct {
				Error string
			}{
				Error: err.Error(),
			})
		} else if newNode, err := sctx.Nodes.AddNode(node); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(struct {
				Error string
			}{
				Error: err.Error(),
			})
		} else {
			ctx.StatusCode(200)
			ctx.JSON(iris.Map{
				"Node": newNode,
				"Message": "Node has been created.",
			})
		}
	})

	// app.Get("/settings", func(ctx iris.Context) {
	// 	if settings, err := sctx.Settings.GetSettings(); err != nil {
	// 		ctx.StatusCode(500)
	// 		ctx.JSON(struct{
	// 			Error string
	// 		}{
	// 			Error: err.Error(),
	// 		})
	// 	} else {
	// 		ctx.JSON(settings)
	// 	}
	// })

	app.Build()
	srv := &http.Server{Handler: app, Addr: addr} // you have to set Handler:app and Addr, see "iris-way" which does this automatically.
	// http://localhost:8080/
	// http://localhost:8080/mypath
	srv.ListenAndServe() // same as app.Run(iris.Addr(":8080"))

	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(addr))
}
