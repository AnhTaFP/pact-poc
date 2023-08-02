package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	type discount struct {
		ID          int      `json:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Value       float64  `json:"value"`
		Vendors     []string `json:"vendors"`
	}

	http.HandleFunc("/discounts/1", func(writer http.ResponseWriter, request *http.Request) {
		d := discount{
			ID:          1,
			Title:       "5.8% off",
			Description: "58th Singaporean National Day discount",
			Type:        "percentage",
			Value:       5.8,
			Vendors:     []string{"vendor1", "vendor2"},
		}
		b, _ := json.Marshal(d)

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(b)
	})

	http.ListenAndServe("localhost:8080", nil)
}
