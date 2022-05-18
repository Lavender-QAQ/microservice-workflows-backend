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

	// build plain/cron workflow depend on starter node
	startNode := (*dag)["Starter"].(*StarterNode)
	ctx := context.Background()
	var err error
	if startNode.WorkflowType == "once" {
		workflow := buildWorkflow(name, templates)
		err = createWorkflow(workflow, ctx, logger)
	} else if startNode.WorkflowType == "cron" {
		cronWorkflow := buildCronWorkflow(name, templates, startNode)
		err = createCronWorkflow(cronWorkflow, ctx, logger)
	}
	if err != nil {
		return err
	}

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

func buildWorkflow(name string, templates []wfv1.Template) wfv1.Workflow {
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
	data, _ := json.Marshal(workflow)
	ioutil.WriteFile("test.json", data, 0666)
	return workflow
}

func createWorkflow(workflow wfv1.Workflow, ctx context.Context, logger logr.Logger) error {
	// create the argo workflow client
	wfClient := WorkflowClient.ArgoprojV1alpha1().Workflows(kubernetes.GetNamespace())
	_, err := wfClient.Create(ctx, &workflow, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Create workflow error")
		return err
	}
	return nil
}

func buildCronWorkflow(name string, templates []wfv1.Template, startNode *StarterNode) wfv1.CronWorkflow {
	var cronExpression string
	var strategy wfv1.ConcurrencyPolicy
	if startNode.CronExpression != "none" {
		cronExpression = startNode.CronExpression
	} else {
		if startNode.TimeUnit == "m" {
			cronExpression = fmt.Sprintf("*/%s * * * *", startNode.Interval)
		} else if startNode.TimeUnit == "h" {
			cronExpression = fmt.Sprintf("* */%s * * *", startNode.Interval)
		} else if startNode.TimeUnit == "d" {
			cronExpression = fmt.Sprintf("* * */%s * *", startNode.Interval)
		}
	}
	if startNode.Strategy == "allow" {
		strategy = wfv1.AllowConcurrent
	} else if startNode.Strategy == "replace" {
		strategy = wfv1.ReplaceConcurrent
	} else if startNode.Strategy == "forbid" {
		strategy = wfv1.ForbidConcurrent
	}
	cronWorkflow := wfv1.CronWorkflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: name,
			Namespace:    kubernetes.GetNamespace(),
		},
		Spec: wfv1.CronWorkflowSpec{
			Schedule:          cronExpression,
			ConcurrencyPolicy: strategy,
			WorkflowSpec: wfv1.WorkflowSpec{
				Entrypoint: name + "-dag",
				Templates:  templates,
			},
		},
	}
	data, _ := json.Marshal(cronWorkflow)
	ioutil.WriteFile("test.json", data, 0666)
	return cronWorkflow
}

func createCronWorkflow(cronWorkflow wfv1.CronWorkflow, ctx context.Context, logger logr.Logger) error {
	// create CronWorkflow
	cronWorkflowClient := WorkflowClient.ArgoprojV1alpha1().CronWorkflows(kubernetes.GetNamespace())
	_, err := cronWorkflowClient.Create(ctx, &cronWorkflow, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Create workflow error")
		return err
	}
	return nil
}
