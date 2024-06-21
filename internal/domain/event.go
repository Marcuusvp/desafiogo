package domain

import (
    "time"
)

type Event struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Organization string    `json:"organization"`
    Date        time.Time `json:"date"`
    Price       int       `json:"price"`
    Rating      string    `json:"rating"`
    ImageURL    string    `json:"image_url"`
    CreatedAt   time.Time `json:"created_at"`
    Location    string    `json:"location"`
}
