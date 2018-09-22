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