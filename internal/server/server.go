package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	E "github.com/IBM/fp-go/either"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))

	return E.Fold(
		func(e error) *http.Server {
			log.Printf("Failed to acquire PORT setting:[%v]\n", e)

			return nil
		},
		func(p int) *http.Server {
			log.Println("Port:", p)

			newServer := &Server{
				port: port,
			}

			// Declare Server config
			server := &http.Server{
				Addr:         fmt.Sprintf(":%d", newServer.port),
				Handler:      newServer.RegisterRoutes(),
				IdleTimeout:  time.Minute,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 30 * time.Second,
			}

			return server
		},
	)(E.TryCatchError(port, err))
}
