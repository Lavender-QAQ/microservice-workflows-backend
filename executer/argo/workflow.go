package argo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/kubernetes"
	"github.com/go-logr/logr"

	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Submit a workflow
func CreateWorkflow(logger logr.Logger, name string, dag *map[string]common.NodeInterface) error {
	logger.Info("Create workflow struct")

	var templates []wfv1.Template
	templates = append(templates, createDAG(logger, name, dag))

	// Write specific nodes in the DAG as templates
	for _, v := range *dag {
		template := v.GenerateTemplate()
		templates = append(templates, template)
	}

	workflow := wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: name,
			Namespace:    kubernetes.GetNamespace(),
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: name + "-dag",
			Templates:  templates,
		},
	}

	// create the argo workflow client
	wfClient := wfclientset.NewForConfigOrDie(kubernetes.GetRestConf()).ArgoprojV1alpha1().Workflows(kubernetes.GetNamespace())
	ctx := context.Background()
	createdWf, err := wfClient.Create(ctx, &workflow, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Create workflow error")
		return err
	}
	fmt.Println(createdWf.Name)

	// // Demo Start
	data, _ := json.Marshal(workflow)
	ioutil.WriteFile("test.json", data, 0666)
	// fmt.Println(string(data))
	// // Demo End

	return nil
}

// Called when creating a workflow to convert the structure of the DAG to argo's DAG type
func createDAG(logger logr.Logger, name string, dag *map[string]common.NodeInterface) v1alpha1.Template {
	logger.Info("Create DAG template")

	var tasks []wfv1.DAGTask

	for _, v := range *dag {
		task := wfv1.DAGTask{
			Name:         v.GetId(),
			Dependencies: v.GetInNode(),
			Template:     v.GetId(),
		}
		if v.HaveInNode() {
			arts := getDAGArtifactsByIncome(v.GetInNode())
			task.Arguments.Artifacts = arts
		}
		tasks = append(tasks, task)
	}

	dags := wfv1.DAGTemplate{
		Tasks: tasks,
	}
	template := wfv1.Template{
		Name: name + "-dag",
		DAG:  &dags,
	}
	return template
}
