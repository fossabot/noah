package api

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/kataras/iris"
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


	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(fmt.Sprintf(":%d", sctx.Flags.HTTPPort)))
}