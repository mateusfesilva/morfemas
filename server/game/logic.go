package game

import (
	"math/rand"
	"strings"
)

var listaMorfemas = []string{"ção", "mente", "dade", "ismo", "ável", "lizar"}

func SortearMorfema() string {
	indice := rand.Intn(len(listaMorfemas))
	return listaMorfemas[indice]
}

func ValidarPalavra(palavra string, morfemaAtual string) bool {
	palavraFormatada := strings.ToLower(strings.TrimSpace(palavra))
	morfemaFormatado := strings.ToLower(strings.TrimSpace(morfemaAtual))

	if len(palavraFormatada) <= len(morfemaFormatado) {
		return false
	}

	if !strings.Contains(palavraFormatada, morfemaFormatado) {
		return false
	}

	return true
}
