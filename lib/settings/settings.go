package settings

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	defaultPort = "7878"
	serverPort = "PORT"
	tgToken = "TOKEN"
	secret = "SECRET"
	tgChannel = "CHANEL_ID"
)

type Settings struct {
	Port 	  string
	TgToken   string
	Secret	  string
	ChannelID int
}

func New(file string) *Settings {
	_ = setEnvFromFile(file)
	port := env(serverPort, defaultPort)
	token := env(tgToken, "")
	secret := env(secret, "")
	channelID, _ := strconv.Atoi(env(tgChannel, ""))
	
	return &Settings{
		Port: 	 port,
		TgToken: token,
		Secret: secret,
		ChannelID: channelID,
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