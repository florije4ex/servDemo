package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type Pong struct {
	gorm.Model
	Message string
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		db, err := gorm.Open("sqlite3", "test.db")
		if err != nil {
			panic("failed to connect database")
		}
		defer func() {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}()

		db.AutoMigrate(&Pong{})
		db.Create(&Pong{Message: "pong"})

		var pong Pong
		db.First(&pong, 1)

		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		res, err := client.Ping().Result()
		log.Println(res)

		c.JSON(200, gin.H{
			"message": pong.Message,
		})
	})
	err := r.Run()
	if err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080
}
