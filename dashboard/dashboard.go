package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DronRathore/goexpress"
	"github.com/go-redis/redis"
)

type Page struct {
	Title string
	Body  []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Hi there, I love %s!", r.URL.RawPath)

}

var redisClient *redis.Client

func main() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	pubsub := redisClient.Subscribe("sensor_data")
	defer pubsub.Close()

	msg := pubsub.Channel()

	go func() {
		//embd.LEDToggle("LED0")
		for {
			//msg, err := redisClient.Get("sensor_data").Result()
			fmt.Println((<-msg).Channel)
			fmt.Println((<-msg).Pattern)
			fmt.Println((<-msg).Payload)
			//fmt.Println(msg.Payload)

			time.Sleep(250 * time.Millisecond)
		}

	}()

	var app = goexpress.Express()

	LibRouter := goexpress.NewRouter()
	app.Use(LibRouter)

	//app.Use(express.Router())

	app.Get("/", func(req goexpress.Request, res goexpress.Response) {
		fmt.Println(req.URL)
		res.SendFile("index.html", true)
	})
	//app.Use(middleware)

	app.Start("8080")

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	//http.HandleFunc("/", handler)

	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)

}
