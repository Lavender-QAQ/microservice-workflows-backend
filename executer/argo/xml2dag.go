package argo

import (
	"errors"
	"sync"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/common"
	"github.com/beevik/etree"
	"github.com/go-logr/logr"
)

var (
	mp       map[string]common.NodeInterface
	mp_mutex sync.Mutex
	doc      *etree.Document
)

// build vertices of DAG
func buildNode(tasks []*etree.Element) {
	var node_wg sync.WaitGroup
	for _, node := range tasks {
		switch node.SelectAttrValue("custom", "none") {
		case "0":
			{
				node_wg.Add(1)
				go buildStarterNode(*node, &node_wg)
				break
			}
		case "1":
			{
				node_wg.Add(1)
				go buildDockerNode(*node, &node_wg)
				break
			}
		case "2":
			{
				node_wg.Add(1)
				go buildPythonscriptNode(*node, &node_wg)
				break
			}
		case "3":
			{
				node_wg.Add(1)
				go buildSuspendNode(*node, &node_wg)
				break
			}
		case "4":
			{
				node_wg.Add(1)
				go buildFilterNode(*node, &node_wg)
				break
			}
		default:
			{
				break
			}
		}
	}
	node_wg.Wait()
}

// Build the edge of the DAG
func buildFlow(flows []*etree.Element) {
	for _, val := range flows {
		sourceRef := val.SelectAttrValue("sourceRef", "none")
		targetRef := val.SelectAttrValue("targetRef", "none")
		// As long as one of the two edges doesn't exist, ignore the edge
		if sourceRef == "none" || targetRef == "none" {
			continue
		}
		mp[sourceRef].AppendOutNode(targetRef)
		mp[targetRef].AppendInNode(sourceRef)
	}
}

// Xml2Dag The entry function of the entire DAG building module,
// completes the initialization of the map,
// loading the etree package, building the DAG, and returning the workflow id
func Xml2Dag(logger logr.Logger, xmlstr string) (string, *map[string]common.NodeInterface, error) {
	mp = make(map[string]common.NodeInterface)
	doc = etree.NewDocument()
	if err := doc.ReadFromString(xmlstr); err != nil {
		logger.Error(err, "Read xmlstr error")
		return "", nil, err
	}
	process := doc.SelectElement("definitions").SelectElement("process")
	processId := process.SelectAttrValue("id", "none")
	buildNode(process.SelectElements("task"))
	buildFlow(process.SelectElements("sequenceFlow"))
	// After DAG is built, map stores the set of all nodes
	if processId == "none" {
		return "", nil, errors.New("Workflow id is not present")
	}
	return processId, &mp, nil
}
