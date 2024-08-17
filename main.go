package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

func SetGINMode() {
	gin.SetMode(gin.DebugMode)
}

func main() {
	SetGINMode()
	router := gin.New()

	router.Use(cors.Default())

	router.Use(apmgin.Middleware(router))

	// Define a simple route
	router.GET("/ping", func(c *gin.Context) {

		apm.DefaultTracer().CaptureHTTPRequestBody(c.Request)
		tx := apm.DefaultTracer().StartTransaction("my-transaction", "custom")

		defer tx.End() 

		ctx := apm.ContextWithTransaction(c, tx)

		performWork(ctx)

		time.Sleep(2 * time.Second)

		panic("[ERROR] JOURNEY SEARCH DUKUH ATAS ke BUNDARAN HI ")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	apmhttp.Wrap(router)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", 8082),
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}


	log.Println("Server running on port 8082")

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func performWork(ctx context.Context) {
	tx := apm.TransactionFromContext(ctx)

	span := tx.StartSpan("performing-work", "custom", nil)
	defer span.End()

	time.Sleep(1 * time.Second)

	fmt.Println("Work has been performed!")
}