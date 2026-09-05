package game

import "fmt"

type Room struct {
	Players map[*Player]bool

	Register   chan *Player
	Unregister chan *Player

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

func (r *Room) Run() {
	for {
		select {
		case player := <-r.Register:
			r.Players[player] = true
			fmt.Println("novo jogador entrou na sala. total: ", len(r.Players))

		case player := <-r.Unregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.Send)
				fmt.Println("um jogador saiu. total: ", len(r.Players))

			}

		case message := <-r.Broadcast:
			for player := range r.Players {
				player.Send <- message
			}
		}

	}
}
