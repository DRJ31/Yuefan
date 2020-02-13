package main

import (
	"github.com/DRJ31/yuefan/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	// Router functions
	controller.Router(router)

	router.Run("127.0.0.1:8003")
}
