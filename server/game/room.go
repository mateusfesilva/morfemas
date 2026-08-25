package game

import "fmt"

// Room centraliza todas as conexões e a lógica do jogo atual
type Room struct {
	// Um mapa (dicionário em Python) de jogadores conectados
	Players map[*Player]bool

	// Canais para sincronizar a entrada e saída de jogadores
	Register   chan *Player
	Unregister chan *Player

	// Canal por onde chegam as palavras enviadas pelos jogadores
	Broadcast chan []byte
}

func NewRoom() *Room {
	return &Room{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Players:    make(map[*Player]bool),
	}
}

// Run é um loop infinito que roda em uma Goroutine separada
func (r *Room) Run() {
	for {
		select {
		case player := <-r.Register:
			r.Players[player] = true
			fmt.Println("Novo jogador conectado! Total:", len(r.Players))

		case player := <-r.Unregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.Send) // Fecha o canal do jogador
				fmt.Println("Jogador saiu. Total:", len(r.Players))
			}

		case message := <-r.Broadcast:
			// Aqui é onde você colocaria a lógica de:
			// 1. Ver de quem veio a palavra
			// 2. Validar no dicionário e checar se tem o morfema
			// 3. Pontuar

			// Por enquanto, apenas retransmite para todos (estilo chat)
			for player := range r.Players {
				player.Send <- message
			}
		}
	}
}
