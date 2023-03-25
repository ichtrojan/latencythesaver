package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		// Initialize the random number generator
		rand.Seed(time.Now().UnixNano())

		// Generate a random integer between 1 and 100
		amount := rand.Intn(100) + 1

		messages := []string{
			"if you've saved this amount before better refresh",
			"garri incomming",
			"better save this",
			"don't cheat",
			"save this one for latency",
			"sugar mummy",
		}

		randIndex := rand.Intn(len(messages))
		message := messages[randIndex]

		response := map[string]string{
			"amount":  fmt.Sprintf("GBP %d", amount),
			"message": message,
		}

		writer.Header().Set("Content-type", "application/json")
		_ = json.NewEncoder(writer).Encode(response)

		return
	})

	if err := http.ListenAndServe(":7777", router); err != nil {
		log.Fatal(err)
	}
}
