package main

import (
	"fmt"
	"log"
	"net/http"
	"ramen/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	redisStore, err := pkg.NewClient("redis:6379", "", 0)
	if err != nil {
		log.Println(err)
	}

	log.Println("Accepting TCP connections")

	r := gin.Default()
	r.GET("/get/:key", func(c *gin.Context) {
		key := c.Param("key")
		message := redisStore.Get(key)
		c.String(http.StatusOK, message)
	})

	r.GET("/set/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		redisStore.Set(key, value)
		message := fmt.Sprintf("Set %s to %s", key, value)
		c.String(http.StatusOK, message)
	})

	r.GET("/del/:key", func(c *gin.Context) {
		key := c.Param("key")
		resp := redisStore.Del(key)
		switch resp {
		case 0:
			c.String(http.StatusOK, fmt.Sprintf("No value for %s", key))
		case 1:
			c.String(http.StatusOK, fmt.Sprintf("Deleted value for %s", key))
		}
	})

	err = r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
	}
}
