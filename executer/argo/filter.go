package argo

import (
	"fmt"
	"sync"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/beevik/etree"
)

type FilterNode struct {
	Node
	leftValue  string
	condition  string
	rightValue string
	expression string
}

func NewFilterNode(id string, leftValue string, condition string, rightValue string, expression string) *FilterNode {
	return &FilterNode{
		Node: Node{
			id:     id,
			custom: 4,
			in:     []string{},
			out:    []string{},
		},
		leftValue:  leftValue,
		condition:  condition,
		rightValue: rightValue,
		expression: expression,
	}
}

// TODO
func (node *FilterNode) GenerateTemplate() v1alpha1.Template {
	parallelStepsList := make([]v1alpha1.ParallelSteps, 0)
	parallelSteps := v1alpha1.ParallelSteps{}
	workflowSteps := make([]v1alpha1.WorkflowStep, 0)
	executeExpression := v1alpha1.WorkflowStep{
		Name:      node.GetId() + "-exec-exp",
		Template:  node.GetId() + "-exec-exp",
		Arguments: v1alpha1.Arguments{},
	}
	judgeExpression := v1alpha1.WorkflowStep{
		Name:      node.GetId() + "-judge-exp",
		Template:  node.GetId() + "-judge-exp",
		Arguments: v1alpha1.Arguments{},
		When:      fmt.Sprintf("{{steps.%s.outputs.result}}", node.GetId()+"-exec-exp"),
	}
	workflowSteps = append(workflowSteps, executeExpression)
	workflowSteps = append(workflowSteps, judgeExpression)
	parallelSteps.Steps = workflowSteps
	parallelStepsList = append(parallelStepsList, parallelSteps)
	template := v1alpha1.Template{
		Name:  node.GetId(),
		Steps: parallelStepsList,
	}

	if node.HaveInNode() && node.HaveOutNode() {
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
	} else if node.HaveOutNode() {
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
	} else if node.HaveInNode() {
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
	}
	return template
}

// Add filter node to map
func buildFilterNode(e etree.Element, node_wg *sync.WaitGroup) {
	defer node_wg.Done()
	id := e.SelectAttrValue("id", "none")
	leftValue := e.SelectAttrValue("leftValue", "none")
	condition := e.SelectAttrValue("condition", "none")
	rightValue := e.SelectAttrValue("rightValue", "none")
	expression := e.SelectAttrValue("expression", "none")
	var node common.NodeInterface = NewFilterNode(id, leftValue, condition, rightValue, expression)
	mp_mutex.Lock()
	mp[id] = node
	mp_mutex.Unlock()
}
