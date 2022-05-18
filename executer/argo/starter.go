package argo

import (
	"fmt"
	"sync"

	v1 "k8s.io/api/core/v1"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/beevik/etree"
)

type StarterNode struct {
	Node
	WorkflowType   string
	Interval       string
	TimeUnit       string
	CronExpression string
	Strategy       string
}

func NewStarterNode(id string, workflowType string, interval string, timeUnit string, cronExpression string,
	stragety string) *StarterNode {
	return &StarterNode{
		Node: Node{
			id:     id,
			custom: 0,
			in:     []string{},
			out:    []string{},
		},
		WorkflowType:   workflowType,
		Interval:       interval,
		TimeUnit:       timeUnit,
		CronExpression: cronExpression,
		Strategy:       stragety,
	}
}

func (node *StarterNode) GenerateTemplate() v1alpha1.Template {
	template := v1alpha1.Template{
		Name: node.GetId(),
		Container: &v1.Container{
			Image: "argoproj/argosay:v2",
			Args: []string{"echo", "{\"msg\":\"workflow starts\"}",
				fmt.Sprintf("/tmp/%s-art.json", node.GetId())},
		},
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

// Add starter node to map
func buildStarterNode(e etree.Element, node_wg *sync.WaitGroup) {
	defer node_wg.Done()
	id := e.SelectAttrValue("id", "none")
	workflowType := e.SelectAttrValue("workflowType", "once")
	interval := e.SelectAttrValue("interval", "none")
	timeUnit := e.SelectAttrValue("timeUnit", "m")
	cronExpression := e.SelectAttrValue("cronExpression", "none")
	strategy := e.SelectAttrValue("strategy", "allow")
	var node common.NodeInterface = NewStarterNode(id, workflowType, interval, timeUnit, cronExpression, strategy)
	mp_mutex.Lock()
	mp[id] = node
	mp_mutex.Unlock()
}
