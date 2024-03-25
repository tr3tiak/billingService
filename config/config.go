package config

import "os"

type Config struct {
	UserDB     string
	PasswordDB string
	NameDB     string
	Port       string
}

func NewConfig() *Config {
	Conf := Config{
		UserDB:     os.Getenv("user_db"),
		PasswordDB: os.Getenv("password_db"),
		NameDB:     os.Getenv("name_db"),
		Port:       "8080",
	}
	return &Conf
}
