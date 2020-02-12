package main

import (
	"github.com/DRJ31/yuefan/controller"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := iris.New()

	// Router functions
	controller.Router(app)

	app.Run(iris.Addr("127.0.0.1:8003"), iris.WithoutServerError(iris.ErrServerClosed))
}
