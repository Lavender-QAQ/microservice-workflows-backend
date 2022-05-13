package argo

import (
	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

// Called when constructing DAGTask to read all incoming edges and generate the corresponding input parameters
func getDAGArtifactsByIncome(income []string) v1alpha1.Artifacts {
	var artifacts v1alpha1.Artifacts
	for _, str := range income {
		artifact := v1alpha1.Artifact{
			Name: str + "-art",
			From: "{{tasks." + str + ".outputs.artifacts." + str + "-art}}",
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts
}

// Called to construct a specific Template for each node,
// iterating over all incoming edges to generate the appropriate input arguments
func getTemplateArtifactsByIncome(income []string) v1alpha1.Artifacts {
	var artifacts v1alpha1.Artifacts
	for _, str := range income {
		artifact := v1alpha1.Artifact{
			Name: str + "-art",
			Path: "/tmp/" + str + "-art",
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts
}

// Called to construct a specific Template for each node, just to determine if there are outsides.
// Return a Artifacts definition based on the id of the node
func getTemplateArtifactsByOutcome(name string) v1alpha1.Artifacts {
	var artifacts v1alpha1.Artifacts
	artifact := v1alpha1.Artifact{
		Name: name + "-art",
		Path: "/tmp/" + name + "-art.json",
	}
	artifacts = append(artifacts, artifact)
	return artifacts
}
