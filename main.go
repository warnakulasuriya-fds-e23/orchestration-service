package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/controller/incomingfingerprintcontroller"
)

func main() {
	_, err := os.Stat(".env")
	if err == nil {
		log.Println("discovered .env file")
		err := godotenv.Load()
		if err != nil {
			log.Println("however failed to load .env file")
		} else {
			log.Println(".env successfully loaded")
		}
	}
	router := gin.Default()

	router.POST("/api/fingerprint/submit-for-identify", incomingfingerprintcontroller.IncomingIdentifyHandler)
	router.POST("/api/fingerprint/submit-for-match", incomingfingerprintcontroller.IncomingMatchHandler)
	router.POST("/api/fingerprint/submit-for-enroll", incomingfingerprintcontroller.IncomingEnrollHandler)

	router.Run(":5000")
}
