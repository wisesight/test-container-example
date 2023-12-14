package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wisesight/test-container-example/cmd/api/handlers"
	"github.com/wisesight/test-container-example/internal/services/post"
	"github.com/wisesight/test-container-example/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	Port         string `env:"PORT" envDefault:"8080"`
	PostMongoURI string `env:"POST_MONGO_URI" envDefault:"mongodb://localhost:27017"`
	PostMongoDB  string `env:"POST_MONGO_DB" envDefault:"posts"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	mongoServerAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.PostMongoURI).SetServerAPIOptions(mongoServerAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("failed to disconnect from mongo: %v", err)
		}
	}()

	postService := post.NewDefaultService(post.NewMongoRepository(client, cfg.PostMongoDB))
	usecases := usecases.NewDefaultUsecases(postService)
	handlers := handlers.NewEchoHandler(usecases)

	server := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	server.GET("/posts/:id", handlers.GetPostByID)

	server.Logger.Fatal(server.Start(":" + cfg.Port))
}
