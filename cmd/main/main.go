package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Nikitapopov/Habbit/internal/config"
	"github.com/Nikitapopov/Habbit/internal/tg"
	"github.com/Nikitapopov/Habbit/internal/user"
	"github.com/Nikitapopov/Habbit/internal/user/db"
	mongodb "github.com/Nikitapopov/Habbit/pkg/client/mongo"
	tg_bot "github.com/Nikitapopov/Habbit/pkg/client/tg"
	"github.com/Nikitapopov/Habbit/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	storagePath     = "storage"
	batchSize       = 100
	tgTokenEnvVar   = "TG_TOKEN"
	tgApiHostEnvVar = "TG_API_HOST"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgMongoDB := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(
		context.Background(),
		logger,
		cfgMongoDB.Host,
		cfgMongoDB.Port,
		cfgMongoDB.Username,
		cfgMongoDB.Password,
		cfgMongoDB.Database,
	)
	if err != nil {
		panic(err)
	}

	logger.Info("register user handler")
	repository := db.NewRepository(mongoDBClient.Collection("users"), logger)
	service := user.NewService(&repository, logger)
	handler := user.NewHandler(&service, logger)
	handler.Register(router)

	tgBot, err := tg_bot.NewBot(cfg.Tg.Token, logger)
	if err != nil {
		panic(err)
	}

	tgClient := tg.NewClient(tgBot, logger)
	tgClient.Start()

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error
	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
