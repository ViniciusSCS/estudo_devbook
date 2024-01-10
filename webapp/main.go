package main

import (
	"fmt"
	"log"
	"net/http"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/generatekey"
	"webapp/src/router"
	"webapp/src/utils"
)

func main() {
	config.Carregar()
	cookies.Configurar()
	utils.CarregarTemplates()

	// Verifica e gera as variáveis de ambiente HASH_KEY e BLOCK_KEY se necessário
	if err := generatekey.GerarKeys(); err != nil {
		log.Fatalf("Erro ao gerar chaves: %v", err)
	}

	r := router.Gerar()

	fmt.Printf("Escutando na porta %d\n", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
