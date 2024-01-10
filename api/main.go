package main

import (
	"api/src/config"
	"api/src/gerakey"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	config.Carregar()

	// Verifica se o comando 'generate' foi fornecido ao executar o programa
	if len(os.Args) > 1 && os.Args[1] == "generate" {
		fmt.Println("\nExecutando 'gerarkey.Gerakey()' para gerar a chave...")
		gerakey.Gerakey()
		fmt.Println("Chave gerada com sucesso.")
		return
	}

	// Se não for o comando 'generate', continua com a inicialização normal
	fmt.Println("Inicializando servidor...")

	// Verifica se a variável de ambiente SECRET_KEY está vazia ou nula
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("SECRET_KEY está vazia ou não definida. Executando 'go run main.go generate'...")
		if err := exec.Command("go", "run", "main.go", "generate").Run(); err != nil {
			log.Fatalf("\nErro ao executar 'go run main.go generate': %v", err)
		}
		fmt.Println("Comando 'go run main.go generate' concluído com sucesso.")
	} else {
		fmt.Printf("SECRET_KEY está definida como %s. Iniciando o servidor...\n", secretKey)
		iniciarServidor()
	}
}

func iniciarServidor() {
	r := router.Gerar()

	fmt.Printf("Escutando na porta %d\n", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
