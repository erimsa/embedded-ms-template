package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	elastic "github.com/olivere/elastic/v6"
	"github.com/rs/cors"
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
		"SensorData":{
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
    }
}`

var sub *redis.Client
var pub *redis.Client

var ctx context.Context

var i int

var client *elastic.Client

func main() {

	r := mux.NewRouter()

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

	var err error
	client, err = GetESClient()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

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

	handler := cors.Default().Handler(r)

	r.HandleFunc("/", createindex)
	//http.HandleFunc("/", handler)
	http.ListenAndServe(":3001", handler)

}

func createindex(w http.ResponseWriter, r *http.Request) {

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
		Type("SensorData").
		BodyJson(sensorData1).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized")

	return client, err

}
