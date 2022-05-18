package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Lavender-QAQ/microservice-workflows-backend/executer"
	"github.com/Lavender-QAQ/microservice-workflows-backend/executer/argo"
	"github.com/go-logr/logr"
)

var HandlerLogger logr.Logger

// Receive post requests sent by the front end
type xmlMsg struct {
	XmlStr string `json:"xml"`
}

// Deployment event triggered by the front-end deployment button
func DeployHandler(w http.ResponseWriter, r *http.Request) {
	logger := HandlerLogger.WithName("DeployHandler")

	// Resolve cross-domain problems
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")

	if r.Method == "POST" {
		logger.Info("Handle post request of deployment")

		var xmlmsg xmlMsg
		data, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err, "Fail to read request body")
			return
		}
		err = json.Unmarshal(data, &xmlmsg)
		if err != nil {
			logger.Error(err, "Fail to parse xml str")
			return
		}

		// Get xml data
		xmlstr := xmlmsg.XmlStr

		// Convert '_' to '-'
		xmlstr = strings.Replace(xmlstr, "_", "-", -1)

		// Parse xml to DAG
		workflowId, mp, err := argo.Xml2Dag(logger.WithName("argo"), xmlstr)
		if err != nil {
			logger.Error(err, "Fail to parse xml")
			return
		}
		// Call executer package
		starter := executer.NewWorkflowStarter(workflowId, mp, logger.WithName("executer"))
		// Create workflow
		err = starter.CreateWorkflow()
		if err != nil {
			logger.Error(err, "Starter create workflow fail")
			return
		}
		_, err = io.WriteString(w, "success")
		if err != nil {
			logger.Error(err, "Fail to write response")
			return
		}
	}
}
