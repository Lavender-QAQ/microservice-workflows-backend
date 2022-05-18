package executer

import (
	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/argo"
	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/go-logr/logr"
)

type WorkflowStarter struct {
	WorkflowId  string
	StarterNode argo.StarterNode
	dag         *map[string]common.NodeInterface
	Logger      logr.Logger
}

// Constructor of the workflow initiator
func NewWorkflowStarter(workflowId string, dag *map[string]common.NodeInterface, logger logr.Logger) *WorkflowStarter {
	return &WorkflowStarter{
		WorkflowId: workflowId,
		dag:        dag,
		Logger:     logger,
	}
}

// Enter the information for the DAG and turn the map into a real workflow
func (w *WorkflowStarter) CreateWorkflow() error {
	logger := w.Logger

	err := argo.CreateWorkflow(w.Logger.WithName("argo"), w.WorkflowId, w.dag)
	if err != nil {
		logger.Error(err, "Create workflow err")
		return err
	}
	return nil
}
