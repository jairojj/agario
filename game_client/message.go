package main

type Message struct {
	ClientID     int `json:"client_id"`
	PlayerCircle PlayerCircle
}
