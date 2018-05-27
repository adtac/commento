package main

import (
	"time"
)

type commenterSession struct {
	Session      string    `json:"session"`
	CommenterHex string    `json:"commenterHex"`
	CreationDate time.Time `json:"creationDate"`
}
