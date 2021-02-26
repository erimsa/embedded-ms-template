package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/googollee/go-socket.io"
)

// SensorData ...
type SensorData struct {
	Service string    `json:"service"` // Service : service adı
	Data    int       `json:"data"`    // Data : servis datası
	Created time.Time `json:"created"` // Created : data nın oluşturulma tarihi
}

const mapping = `{
    "settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
    "mappings":{
        "properties":{
                "service":{
                    "type": "text"
                },
                "data":{
                    "type": "integer"
                },
                "created":{
                    "type": "date"
                }
        }
    }
}`

var sub *redis.Client
var pub *redis.Client

var client *elastic.Client

var ctx context.Context

var i int

func main() {

	ctx = context.Background()

	/*
		server, err := socketio.NewServer(nil)
		if err != nil {
			log.Fatal(err)
		}

		/*
			server.On("connection", func(so socketio.Socket) {
				so.Join("sensor_data")
				so.On("message", func(msg string) {
					so.Emit("message", msg)
					so.BroadcastTo("sensor_data", "", args)
				})
			})
	*/

	sub = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pub = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if _, err := sub.Ping().Result(); err != nil {
		panic(err)
		//fmt.Println(subPong, err)
		//os.Exit(0)
	}

	if _, err := pub.Ping().Result(); err != nil {
		panic(err)
		//fmt.Println(pubPong, err)
		//os.Exit(0)
	}

	{
		var err error
		client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://elasticsearch:9200"))
		if err != nil {
			panic(err)
		}
	}
	info, code, err := client.Ping("http://elasticsearch:9200").Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("eliar").Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		createIndex, err := client.CreateIndex("eliar").BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}

		if !createIndex.Acknowledged {
			fmt.Println("salih")
		}
	}

	_, err = client.Flush().Index("eliar").Do(ctx)
	if err != nil {
		panic(err)
	}

	subsub := sub.Subscribe("sensor_data")
	defer subsub.Close()

	/*
		for i := 0; i < 10000000; i++ {
			pub.Publish("sensor_data", strconv.Itoa(i))

		}
	*/

	http.HandleFunc("/", handler)
	//http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {

	i++

	rdata := rand.Intn(100)
	tTime := time.Now()
	//data := strconv.Itoa(i)
	sensorData1 := &SensorData{Service: "heat-sensor-service", Data: rdata, Created: tTime}

	bytes, err := json.Marshal(sensorData1)
	if err != nil {
		panic(err)
	}
	s := string(bytes[:])

	pub.Publish("sensor_data", s)

	put1, err := client.Index().
		Index("eliar").
		Type("_doc").
		BodyJson(sensorData1).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
