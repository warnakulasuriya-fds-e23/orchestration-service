package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/config"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/controller/incomingfingerprintcontroller"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/controller/outgoingfingerprintcontroller"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/customstorage"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		log.Println(string(body))
		log.Println(c.Request.Header)
		c.Next()
	}
}

func testHandler(c *gin.Context) {
	log.Println("recieved!")
	c.IndentedJSON(http.StatusOK, "recievd")
}

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
	devicesConfigFilePath := os.Getenv("DEVICES_CONFIG_JSON_PATH")
	if devicesConfigFilePath == "" {
		log.Fatalf("devices config json path not specified in environment variable DEVICES_CONFIG_JSON_PATH")
	}

	devicesConfig, errDevConfigLoader := config.DeviceConfigLoader(devicesConfigFilePath)
	if errDevConfigLoader != nil {
		log.Fatalf("error while trying to load up devices config json : %s", err.Error())
	}

	tokenstorage, err := customstorage.NewTokenStorage()
	if err != nil {
		log.Fatalf("unable to make token storage : %s", err.Error())
	}
	outgoingfingerprintcntrlr := outgoingfingerprintcontroller.NewOutgoingFingerprintController(devicesConfig, *tokenstorage)
	incomingfingerprintcntrlr := incomingfingerprintcontroller.NewIncomingFingerprintController(devicesConfig, outgoingfingerprintcntrlr)

	router := gin.Default()
	router.Use(RequestLoggerMiddleware())
	router.POST("/api/fingerprint/identify", incomingfingerprintcntrlr.IncomingIdentifyHandler)
	router.POST("/api/fingerprint/match", incomingfingerprintcntrlr.IncomingMatchHandler)
	router.POST("/api/fingerprint/enroll", incomingfingerprintcntrlr.IncomingEnrollHandler)
	router.POST("/api/fingerprint/authorize", incomingfingerprintcntrlr.IncomingAuthorize)

	router.POST("/oauth2/authorize", testHandler)
	router.Run(":5000")
}
