package argo

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/beevik/etree"
	v1 "k8s.io/api/core/v1"
)

const RequestImage string = "lavenderqaq/agency"

// Docker node
type DockerNode struct {
	Node
	Image   string
	Port    int
	Target  string
	Command []string
	Args    []string
}

func NewDockerNode(id string, image string, port int, target string, command []string, args []string) *DockerNode {
	return &DockerNode{
		Node: Node{
			id:     id,
			custom: 1,
			in:     []string{},
			out:    []string{},
		},
		Image:   image,
		Port:    port,
		Target:  target,
		Command: command,
	}
}

func (node *DockerNode) GenerateTemplate() v1alpha1.Template {
	template := v1alpha1.Template{
		Name: node.GetId(),
		Container: &v1.Container{
			Image:   RequestImage,
			Command: []string{"./agency"},
		},
	}

	var args Args
	if node.HaveInNode() && node.HaveOutNode() {
		args = NewArgsWithInputAndOutput(node.GetInNode()[0], node.GetId(), node.Image, node.Port, node.Target)
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
	} else if node.HaveOutNode() {
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
		args = NewArgsWithOutput(node.GetId(), node.Image, node.Port, node.Target)
	} else if node.HaveInNode() {
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
		args = NewArgsWithInput(node.GetInNode()[0], node.Image, node.Port, node.Target)
	}
	template.Container.Args = args

	return template
}

func (node *DockerNode) GenerateTemplate() v1alpha1.Template {
	template := v1alpha1.Template{
		Name: node.GetId(),
		Container: &v1.Container{
			Image:   RequestImage,
			Command: []string{"./agency"},
		},
	}

	var args Args
	if node.HaveInNode() && node.HaveOutNode() {
		args = NewArgsWithInputAndOutput(node.GetInNode()[0], node.GetId(), node.Image, node.Port, node.Target)
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
	} else if node.HaveOutNode() {
		template.Outputs.Artifacts = getTemplateArtifactsByOutcome(node.GetId())
		args = NewArgsWithOutput(node.GetId(), node.Image, node.Port, node.Target)
	} else if node.HaveInNode() {
		template.Inputs.Artifacts = getTemplateArtifactsByIncome(node.GetInNode())
		args = NewArgsWithInput(node.GetInNode()[0], node.Image, node.Port, node.Target)
	}
	template.Container.Args = args

	return template
}

// Add docker node to map
func buildDockerNode(e etree.Element, node_wg *sync.WaitGroup) {
	defer node_wg.Done()
	id := e.SelectAttrValue("id", "none")
	image := e.SelectAttrValue("image", "none")
	portStr := e.SelectAttrValue("port", "none")

	target := e.SelectAttrValue("target", "/")
	commandStr := e.SelectAttrValue("command", "none")
	argsStr := e.SelectAttrValue("args", "none")

	image = "httpserver"

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 80
	}
	command := strings.Split(commandStr, " ")
	args := strings.Split(argsStr, " ")

	var node common.NodeInterface = NewDockerNode(id, image, port, target, command, args)
	mp_mutex.Lock()
	mp[id] = node
	mp_mutex.Unlock()
}
