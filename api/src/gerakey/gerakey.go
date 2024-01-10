package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// EnvFileWriter é responsável por escrever no arquivo .env.
type EnvFileWriter interface {
	WriteToEnvFile(filePath, key, value string) error
}

// EnvFileReader é responsável por ler o arquivo .env.
type EnvFileReader interface {
	ReadEnvFile(filePath, key string) (string, error)
}

// EnvFileWriterImpl implementa EnvFileWriter.
type EnvFileWriterImpl struct{}

// WriteToEnvFile escreve no arquivo .env.
func (w *EnvFileWriterImpl) WriteToEnvFile(filePath, key, value string) error {
	// Lê o arquivo .env existente (se houver)
	err := godotenv.Load(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Abre o arquivo .env para escrever, mas não cria se não existir
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	// Adiciona ou atualiza a chave no mapa
	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// Verifica se a chave já existe
	if _, exists := envMap[key]; exists {
		// Atualiza o valor da chave
		envMap[key] = value
	} else {
		// Adiciona a nova chave
		envMap[key] = value
	}

	// Volta ao início do arquivo
	file.Seek(0, 0)
	file.Truncate(0)

	// Escreve o mapa no arquivo .env
	for k, v := range envMap {
		fmt.Fprintf(file, "%s=%s\n", k, v)
	}

	return nil
}

// EnvFileReaderImpl implementa EnvFileReader.
type EnvFileReaderImpl struct{}

// ReadEnvFile lê o valor de uma chave no arquivo .env.
func (r *EnvFileReaderImpl) ReadEnvFile(filePath, key string) (string, error) {
	// Lê o arquivo .env existente (se houver)
	err := godotenv.Load(filePath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// Procura o valor da chave
	for _, env := range os.Environ() {
		kv := splitEnv(env)
		if kv[0] == key {
			return kv[1], nil
		}
	}

	return "", nil
}

func splitEnv(env string) []string {
	return strings.SplitN(env, "=", 2)
}

func main() {
	app := cli.NewApp()
	app.Name = "Gerador de SECRET_KEY"
	app.Usage = "Gera uma SECRET_KEY em base64 e a salva ou atualiza no arquivo .env"

	// Inicializa os implementadores
	envWriter := &EnvFileWriterImpl{}
	envReader := &EnvFileReaderImpl{}

	app.Commands = []cli.Command{
		{
			Name:   "generate",
			Usage:  "Gera e salva ou atualiza a SECRET_KEY",
			Action: GenerateOrUpdateSecretKey(envWriter, envReader),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

// GenerateOrUpdateSecretKey é uma função que encapsula a lógica para gerar e adicionar ou atualizar a SECRET_KEY ao .env.
func GenerateOrUpdateSecretKey(writer EnvFileWriter, reader EnvFileReader) cli.ActionFunc {
	return func(c *cli.Context) error {
		// Gera uma chave aleatória
		secretKey, err := generateRandomKey()
		if err != nil {
			log.Fatal(err)
		}

		// Codifica a chave em base64
		encodedKey := encodeToString(secretKey)

		// Verifica se a chave já existe no .env
		currentKey, err := reader.ReadEnvFile("../../.env", "SECRET_KEY")
		if err != nil {
			return err
		}

		// Se a chave já existir, atualiza o valor
		if currentKey != "" {
			// Atualiza a chave no arquivo .env
			err := writer.WriteToEnvFile("../../.env", "SECRET_KEY", `"`+encodedKey+`"`)
			if err != nil {
				return err
			}
			fmt.Println("SECRET_KEY atualizada com sucesso no arquivo .env!")
		} else {
			// Se a chave não existir, adiciona ao arquivo .env
			err := writer.WriteToEnvFile("../../.env", "SECRET_KEY", encodedKey)
			if err != nil {
				return err
			}
			fmt.Println("SECRET_KEY gerada e salva com sucesso no arquivo .env!")
		}

		return nil
	}
}

// generateRandomKey gera uma chave aleatória.
func generateRandomKey() ([]byte, error) {
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// encodeToString codifica dados em base64.
func encodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
