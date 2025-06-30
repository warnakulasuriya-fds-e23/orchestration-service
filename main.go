package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/controller/incomingfingerprintcontroller"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/controller/outgoingfingerprintcontroller"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/customstorage"
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
	tokenstorage, err := customstorage.NewTokenStorage()
	if err != nil {
		log.Fatalf("unable to make token storage : %s", err.Error())
	}
	outgoingfingerprintcntrlr := outgoingfingerprintcontroller.NewOutgoingFingerprintController(*tokenstorage)
	incomingfingerprintcntrlr := incomingfingerprintcontroller.NewIncomingFingerprintController(outgoingfingerprintcntrlr)

	router := gin.Default()

	router.POST("/api/fingerprint/identify", incomingfingerprintcntrlr.IncomingIdentifyHandler)
	router.POST("/api/fingerprint/match", incomingfingerprintcntrlr.IncomingMatchHandler)
	router.POST("/api/fingerprint/enroll", incomingfingerprintcntrlr.IncomingEnrollHandler)

	router.Run(":5000")
}
