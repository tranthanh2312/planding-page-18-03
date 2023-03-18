package main

import (
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/Tech-by-GL/dashboard/db"
	"github.com/Tech-by-GL/dashboard/handler"
	"github.com/Tech-by-GL/dashboard/middleware"
	"github.com/Tech-by-GL/dashboard/singleton"
)

func init() {
	color.Red("Starting server at port 8082...")
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	singleton.InitTimeLocation()

	r.Use(middleware.JSONWriterMiddleware)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	handler := handler.NewHandler()
	db.InitBoltDB()

	r.GET("/tuition", handler.GetTuition)
	r.GET("/invoice", handler.GetKPIInvoice)

	// ! Warning system
	r.GET("/warning/orders", handler.GetIncompatibleOrder)
	r.GET("/warning/orders_future", handler.GetWarningOrderWithFutureDate)
	r.GET("/warning/duplicated_order", handler.GetWarningDuplicatedOrders)

	err := r.Run(":8082")
	if err != nil {
		panic(err)
	}
}
