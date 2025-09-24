package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Komilov31/event-booker/internal/config"
	"github.com/Komilov31/event-booker/internal/handler"
	"github.com/Komilov31/event-booker/internal/rabbitmq"
	"github.com/Komilov31/event-booker/internal/repository"
	"github.com/Komilov31/event-booker/internal/sender"
	"github.com/Komilov31/event-booker/internal/service"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	zlog.Init()

	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.Postgres.Host,
		config.Cfg.Postgres.Port,
		config.Cfg.Postgres.User,
		config.Cfg.Postgres.Password,
		config.Cfg.Postgres.Name,
	)

	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		log.Fatal("could not init db: " + err.Error())
	}

	repository := repository.New(db)
	queue := rabbitmq.New()
	sender := sender.New()
	service := service.New(repository, queue, sender)
	handler := handler.New(service)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting signal %v. Shuting down", sig)
		cancel()
	}()

	service.StartWorker(ctx)

	router := ginext.New()
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	// Register static files
	engine.LoadHTMLFiles("/app/static/index.html", "/app/static/admin.html", "/app/static/user.html")
	engine.Static("/static", "/app/static")

	// POST requests
	engine.POST("/events", handler.CreateEvent)
	engine.POST("/events/:id/book", handler.CreateBooking)
	engine.POST("/events/:id/confirm", handler.ConfirmPayment)

	// GET requests

	// Front pages
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/admin", handler.GetAdminPage)
	engine.GET("/user", handler.GetUserPage)

	// API endpoints
	engine.GET("/events/:id", handler.GetEvent)
	engine.GET("/events", handler.GetAllEvents)

}
