package argo

import (
	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/kubernetes"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
)

var WorkflowClient *wfclientset.Clientset

func Init() {
	WorkflowClient = wfclientset.NewForConfigOrDie(kubernetes.GetRestConf())
}
