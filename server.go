package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	database "github.com/ichtrojan/latencythesaver/redis"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Amount  int
	Message string
}

func main() {
	router := chi.NewRouter()

	_ = godotenv.Load()

	redisHost, exist := os.LookupEnv("REDIS_HOST")

	if !exist {
		log.Fatal("REDIS_HOST not set in .env")
	}

	redisPort, exist := os.LookupEnv("REDIS_PORT")

	if !exist {
		log.Fatal("REDIS_PORT not set in .env")
	}

	redisPass, exist := os.LookupEnv("REDIS_PASS")

	if !exist {
		log.Fatal("REDIS_PASS not set in .env")
	}

	redisScheme, exist := os.LookupEnv("REDIS_SCHEME")

	if !exist {
		log.Fatal("REDIS_SCHEME not set in .env")
	}

	if err := database.ConnectRedis(redisHost, redisPort, redisPass, redisScheme); err != nil {
		log.Fatal(err)
	}

	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		var amount int

		var message string

		var savedAmounts []int

		latency, _ := database.Redis.Get("latency_the_saver").Result()

		_ = json.Unmarshal([]byte(latency), &savedAmounts)

		for {
			// Initialize the random number generator
			rand.Seed(time.Now().UnixNano())

			// Generate a random integer between 1 and 100
			amount = rand.Intn(100) + 1

			messages := []string{
				"if you've saved this amount before better refresh",
				"garri incomming",
				"better save this",
				"don't cheat",
				"save this one for latency",
				"sugar mummy",
			}

			randIndex := rand.Intn(len(messages))
			message = messages[randIndex]

			if contains(savedAmounts, amount) == false {
				break
			}
		}

		tmpl := template.Must(template.ParseFiles("index.html"))

		data := Data{
			Amount:  amount,
			Message: message,
		}

		_ = tmpl.Execute(writer, data)
	})

	router.Post("/save", func(writer http.ResponseWriter, request *http.Request) {
		amount, _ := strconv.Atoi(request.FormValue("amount"))

		latency, _ := database.Redis.Get("latency_the_saver").Result()

		if latency == "" {
			data := []int{amount}

			jsonData, _ := json.Marshal(data)

			database.Redis.Set("latency_the_saver", string(jsonData), 0)
		} else {
			var savedAmounts []int

			_ = json.Unmarshal([]byte(latency), &savedAmounts)

			savedAmounts = append(savedAmounts, amount)

			jsonData, _ := json.Marshal(savedAmounts)

			err := database.Redis.Set("latency_the_saver", string(jsonData), 0).Err()

			if err != nil {
				fmt.Println(err)
			}
		}

		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
	})

	if err := http.ListenAndServe(":7777", router); err != nil {
		log.Fatal(err)
	}
}

func contains(arr []int, searchValue int) bool {
	found := false

	for _, value := range arr {
		if value == searchValue {
			found = true
			break
		}
	}

	return found
}
