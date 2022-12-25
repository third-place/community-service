/*
 * Community service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/third-place/community-service/internal"
	"github.com/third-place/community-service/internal/middleware"
	"github.com/third-place/community-service/internal/service"
	"log"
	"net/http"
)

func main() {
	go readKafka()
	serveHttp()
}

func readKafka() {
	log.Print("connecting to kafka")
	service.CreateConsumerService().Loop()
	log.Print("exit kafka loop")
}

func serveHttp() {
	router := internal.NewRouter()
	handler := cors.AllowAll().Handler(router)
	log.Print("http listening on 8081")
	log.Fatal(http.ListenAndServe(":8081",
		middleware.ContentTypeMiddleware(handler)))
}
