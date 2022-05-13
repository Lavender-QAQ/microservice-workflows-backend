package argo

import (
	"fmt"
	"os"
	"sync"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/beevik/etree"
	v1 "k8s.io/api/core/v1"
)

type PythonscriptNode struct {
	Node
	version string
	script  string
}

func NewPythonscriptNode(id string, version string, script string) *PythonscriptNode {
	return &PythonscriptNode{
		Node: Node{
			id:     id,
			custom: 2,
			in:     []string{},
			out:    []string{},
		},
		version: version,
		script:  script,
	}
}

func (node *PythonscriptNode) GenerateTemplate() v1alpha1.Template {
	harborUrl := os.Getenv("HARBOR_URL")
	image := fmt.Sprintf("%s/argo/python:%s", harborUrl, node.version)
	// add prefix of reading json file and suffix of writing json file
	script := addScriptPrefixAndSuffix(node)
	template := v1alpha1.Template{
		Name: node.GetId(),
		Script: &v1alpha1.ScriptTemplate{
			Container: v1.Container{
				Image:   image,
				Command: []string{"python"},
			},
			Source: script,
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

// Add pythonscript node to map
func buildPythonscriptNode(e etree.Element, node_wg *sync.WaitGroup) {
	defer node_wg.Done()
	id := e.SelectAttrValue("id", "none")
	version := e.SelectAttrValue("version", "3.6")
	script := e.SelectAttrValue("script", "none")
	var node common.NodeInterface = NewPythonscriptNode(id, version, script)
	mp_mutex.Lock()
	mp[id] = node
	mp_mutex.Unlock()
}

func addScriptPrefixAndSuffix(node *PythonscriptNode) string {
	script := node.script
	prefix := "import json\n\n"
	if len(node.GetInNode()) == 1 {
		inputJsonPath := fmt.Sprintf("/tmp/%s-art", node.GetInNode()[0])
		prefix += fmt.Sprintf("input = json.load(open('%s', 'r', encoding='utf-8'))\n\n", inputJsonPath)
	}

	outputJsonPath := fmt.Sprintf("/tmp/%s-art.json", node.GetId())
	suffix := fmt.Sprintf("\n\ntry:\n"+
		"    json.dump(result, open('%s', 'w'))\n"+
		"    print(result)\n"+
		"    print('success')\n"+
		"except Exception as e:\n"+
		"    print(e)", outputJsonPath)
	script = prefix + script + suffix

	return script
}
