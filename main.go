package main

import (
	"fmt"
	"log"
	"net/http"
	"ramen/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	redisStore, err := pkg.NewClient("redis:6379", "", 0)
	if err != nil {
		log.Println(err)
	}

	log.Println("Accepting TCP connections")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ENDPOINTS:\nGET: /get/:key\nSET: /set/:key/:value\nDEL /del/:key")
		// c.String(http.StatusOK, "alive")
	})

	r.GET("/get/:key", func(c *gin.Context) {
		key := c.Param("key")
		message, err := redisStore.Get(key)
		if err == nil {
			c.String(http.StatusOK, message)
		} else {
			log.Println(err)
			c.String(http.StatusServiceUnavailable, "Service Unavailable")
		}
	})

	r.GET("/set/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		err := redisStore.Set(key, value)
		if err == nil {
			message := fmt.Sprintf("Set %s to %s", key, value)
			c.String(http.StatusOK, message)
		} else {
			c.String(http.StatusServiceUnavailable, "Service Unavailable")
		}
	})

	r.GET("/del/:key", func(c *gin.Context) {
		key := c.Param("key")
		resp, err := redisStore.Del(key)
		if err == nil {
			switch resp {
			case 0:
				c.String(http.StatusOK, fmt.Sprintf("No value for %s", key))
			case 1:
				c.String(http.StatusOK, fmt.Sprintf("Deleted value for %s", key))
			}
		} else {
			c.String(http.StatusServiceUnavailable, "Service Unavailable")
		}
	})

	r.GET("/slow/get/:key", func(c *gin.Context) {
		key := c.Param("key")
		time.Sleep(5 * time.Second)
		message, err := redisStore.Get(key)
		if err == nil {
			c.String(http.StatusOK, message)
		} else {
			c.String(http.StatusServiceUnavailable, "Service Unavailable")
		}
	})

	err = r.Run(fmt.Sprintf("0.0.0.0:8080"))
	if err != nil {
		log.Fatalln(err)
	}
}
