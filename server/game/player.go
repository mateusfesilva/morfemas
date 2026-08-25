package game

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tempo máximo que o servidor espera para conseguir enviar uma mensagem
	writeWait = 10 * time.Second
	// Tempo máximo esperando um sinal de vida (pong) do frontend
	pongWait = 60 * time.Second
	// Frequência que o servidor manda um "ping" para ver se o jogador ainda tá lá
	pingPeriod = (pongWait * 9) / 10
	// Tamanho máximo em bytes de uma palavra que o jogador pode enviar
	maxMessageSize = 512
)

// Player representa a conexão do usuário e seus canais de comunicação
type Player struct {
	Room *Room
	Conn *websocket.Conn
	// O tubo por onde as mensagens do servidor vão chegar até esse jogador
	Send chan []byte
}

// NewPlayer inicializa o jogador (Essa é a função que o Go estava pedindo!)
func NewPlayer(room *Room, conn *websocket.Conn) *Player {
	return &Player{
		Room: room,
		Conn: conn,
		Send: make(chan []byte, 256), // Buffer de 256 mensagens para evitar engasgos
	}
}

// ReadPump fica escutando tudo o que o frontend do jogador enviar (ex: as palavras)
func (p *Player) ReadPump() {
	// O defer garante que, quando o jogador fechar a aba, ele sai da sala e a conexão encerra
	defer func() {
		p.Room.Unregister <- p
		p.Conn.Close()
	}()

	p.Conn.SetReadLimit(maxMessageSize)
	p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	p.Conn.SetPongHandler(func(string) error { p.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Loop infinito escutando as mensagens
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro de fechamento inesperado: %v", err)
			}
			break
		}

		// Limpa a mensagem (tira espaços vazios das pontas) e manda para o tubo principal da Sala
		message = bytes.TrimSpace(message)
		p.Room.Broadcast <- message
	}
}

// WritePump pega as mensagens da Sala e envia para a tela (frontend) do jogador
func (p *Player) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.Send:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Se o canal "Send" foi fechado, manda o frontend desconectar
				p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message) // Envia a mensagem pro navegador!

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Envia um "ping" silencioso pro navegador pra manter a conexão viva
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
