package main

import (
	"fmt"
	"log"
	"net/http"
	"server/game"

	"github.com/gorilla/websocket"
)

// Upgrader transforma a conexão HTTP comum em um WebSocket em tempo real
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permite conexões de qualquer frontend (útil para desenvolvimento)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(room *game.Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao criar websocket:", err)
		return
	}

	// Cria um novo jogador associado a esta conexão e à sala
	player := game.NewPlayer(room, conn)

	// Registra o jogador na sala
	room.Register <- player

	// Inicia as Goroutines que ficam escutando e enviando mensagens
	// O comando "go" faz elas rodarem em paralelo sem travar o servidor
	go player.ReadPump()
	go player.WritePump()
}

func main() {
	// Cria a sala e coloca ela para rodar em background
	room := game.NewRoom()
	go room.Run()

	// Define a rota do WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(room, w, r)
	})

	fmt.Println("Servidor de WebSockets rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
