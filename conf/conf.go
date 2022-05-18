package conf

import (
	"os"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/argo"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/kubernetes"
	"github.com/joho/godotenv"
	ctrl "sigs.k8s.io/controller-runtime"
)

var logger = ctrl.Log.WithName("")

func Init(namespace string) {
	// 从本地读取环境变量
	_ = godotenv.Load()
	if os.Getenv("ACTIVE_ENV") == "DEV" {
		_ = godotenv.Load(".env.dev")
	} else if os.Getenv("ACTIVE_ENV") == "PROD" {
		_ = godotenv.Load(".env.prod")
	}

	// Init client-go
	err := kubernetes.Init(namespace)
	if err != nil {
		logger.Error(err, "Fail to initialize kubernetes cluster")
		return
	}

	// Init argo-workflow
	argo.Init()
}
