package service

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadDotEnv(filenames ...string) {
	err1 := godotenv.Load(filenames...)
	if err1 != nil {
		log.Fatalln("Error loading .env file")
	}
}
