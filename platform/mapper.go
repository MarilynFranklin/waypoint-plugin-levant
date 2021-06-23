package platform

import (
	"github.com/hashicorp/waypoint/builtin/nomad"
)

// LevantDeploymentMapper maps a Deployment to a nomad.Deployment.
func LevantDeploymentMapper(src *Deployment) *nomad.Deployment {
	return &nomad.Deployment{
		Id:   src.Id,
		Name: src.Name,
	}
}
