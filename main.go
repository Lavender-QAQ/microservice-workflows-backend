package main

import (
	"flag"
	"fmt"

	"github.com/Lavender-QAQ/microservice-workflows-backend/handler"
	"github.com/Lavender-QAQ/microservice-workflows-backend/router"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger logr.Logger

func main() {
	kubeconfigPath := flag.String("kubeconfig", "./kubeconfig", "Kubernetes configuration file location")

	flag.Parse()

	err := registerLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.WithValues("kubeconfig location", *kubeconfigPath).Info("Kubeconfig parameters were successfully parsed")

	err = router.NewRouter("127.0.0.1:30086")
	if err != nil {
		logger.Error(err, "Fail to create router")
		return
	}
}

func registerLogger() error {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("who watches the watchmen (%v)?", err)
	}
	logger = zapr.NewLogger(zapLog)

	// Register handler
	handler.HandlerLogger = logger.WithName("Handler")
	router.RouterLogger = logger.WithName("Router")

	return nil
}
