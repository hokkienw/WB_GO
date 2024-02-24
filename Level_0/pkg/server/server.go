package server

import (
	"net/http"
	"level_0/pkg/cache"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
)

func RunServer(cache_instance *cache.Cache) {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html");

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/:orderID", func(c *gin.Context) {
		orderID := c.Param("orderID")
	
		data, found := cache_instance.Get(string(orderID))
		if !found {
			log.Printf("Order not found for orderID: %s", data)
			c.HTML(http.StatusNotFound, "error-template.html", gin.H{"Message": "Order not found"})
			return
		}
		totalOrders := cache_instance.Count()
		pageNum,_ := strconv.Atoi(orderID)

		c.HTML(http.StatusOK, "template.html", gin.H{
			"Total":     "Total: " + strconv.Itoa(totalOrders),
			"OrderData": data.(string),
			"PageNum":   pageNum,
			"TotalPage": totalOrders,
			"PrevPage":  max(pageNum-1, 1),
			"NextPage":  min(pageNum+1, totalOrders),
			"LastPage":  totalOrders,
		})
	})

	router.Run(":8080")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}