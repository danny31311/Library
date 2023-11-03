package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"library/api"
	"library/api/middleware"
	"library/cfg"
	"library/internals/app/db"
	"library/internals/app/handlers"
	"library/internals/app/processors"
	"log"
	"net/http"
	"time"
)

type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.config = config
	server.ctx = ctx
	return server
}

func (server *Server) Serve() {
	log.Println("Starting server")
	var err error
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBString())
	if err != nil {
		log.Fatal(err)
	}

	booksStorage := db.NewBooksStorage(server.db)
	usersStorage := db.NewUsersStorage(server.db)

	booksProcessor := processors.NewBooksProcessor(booksStorage)
	usersProcessor := processors.NewUsersProcessor(usersStorage)

	userHandler := handlers.NewUsersHandler(usersProcessor)
	booksHandler := handlers.NewBooksHandler(booksProcessor)

	routes := api.CreateRoutes(userHandler, booksHandler)
	routes.Use(middleware.RequestLog)

	server.srv = &http.Server{Addr: ":" + server.config.Port, Handler: routes}
	log.Println("Serve started")
	err = server.srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) ShutDown() {
	log.Println("Serve stopped")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Second*5)
	server.db.Close()
	defer func() {
		cancel()
	}()
	var err error
	if err = server.srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Serve shutdown failed: %v", err)
	}
	log.Println("server exited properly")

}
