package settings

import (
	"bufio"
	"os"
	"strings"
)

const (
	defaultPort = "7878"
	serverPort = "PORT"
	tgToken = "TOKEN"
)

type Settings struct {
	Port 	 string
	TgToken string
}

func New(file string) *Settings {
	_ = setEnvFromFile(file)
	port := env(serverPort, defaultPort)
	token := env(tgToken, "")
	
	return &Settings{
		Port: 	 port,
		TgToken: token,
	}
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func setEnvFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}