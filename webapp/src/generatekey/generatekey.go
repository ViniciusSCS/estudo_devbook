package generatekey

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
)

// GerarKeys verifica e gera as variáveis de ambiente HASH_KEY e BLOCK_KEY se não estiverem definidas.
func GerarKeys() error {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Verifica e gera a variável de ambiente HASH_KEY se necessário
	hashKey := os.Getenv("HASH_KEY")
	if hashKey == "" {
		// Se HASH_KEY não estiver definida, gera a chave
		fmt.Println("HASH_KEY não está definida. Gerando chave...")
		hashKey, err := gerarChave()
		if err != nil {
			return fmt.Errorf("Erro ao gerar a chave HASH_KEY: %v", err)
		}
		fmt.Printf("HASH_KEY gerada e definida como %s\n", hashKey)

		// Atualiza a variável de ambiente no arquivo .env
		if err := writeEnvVariable("HASH_KEY", hashKey); err != nil {
			return fmt.Errorf("Erro ao escrever no arquivo .env: %v", err)
		}
	} else {
		fmt.Printf("HASH_KEY já está definida como %s. Não será gerada novamente.\n", hashKey)
	}

	// Verifica e gera a variável de ambiente BLOCK_KEY se necessário
	blockKey := os.Getenv("BLOCK_KEY")
	if blockKey == "" {
		// Se BLOCK_KEY não estiver definida, gera a chave
		fmt.Println("BLOCK_KEY não está definida. Gerando chave...")
		blockKey, err := gerarChave()
		if err != nil {
			return fmt.Errorf("Erro ao gerar a chave BLOCK_KEY: %v", err)
		}
		fmt.Printf("BLOCK_KEY gerada e definida como %s\n", blockKey)

		// Atualiza a variável de ambiente no arquivo .env
		if err := writeEnvVariable("BLOCK_KEY", blockKey); err != nil {
			return fmt.Errorf("Erro ao escrever no arquivo .env: %v", err)
		}
	} else {
		fmt.Printf("BLOCK_KEY já está definida como %s. Não será gerada novamente.\n", blockKey)
	}

	return nil
}

// gerarChave gera uma chave aleatória.
func gerarChave() (string, error) {
	chave := securecookie.GenerateRandomKey(16)
	_, err := rand.Read(chave)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(chave), nil
}

// writeEnvVariable escreve a variável de ambiente no arquivo .env
func writeEnvVariable(key, value string) error {
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		return err
	}

	return nil
}
