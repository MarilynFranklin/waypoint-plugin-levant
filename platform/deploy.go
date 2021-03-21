package platform

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/levant/helper"
	"github.com/hashicorp/levant/levant"
	"github.com/hashicorp/levant/levant/structs"
	"github.com/hashicorp/levant/template"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const (
	metaId    = "waypoint.hashicorp.com/id"
	metaNonce = "waypoint.hashicorp.com/nonce"
)

type DeployConfig struct {
	// The HTTP API endpoint for Nomad where all calls will be made
	// ("http://localhost:4646").
	Address string `hcl:"address,optional"`

	// The Consul host and port to use when making Consul KeyValue lookups
	// for template rendering ("localhost:8500").
	ConsulAddress string `hcl:"consul_address,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Variables that are meant to configure the nomad job template file.
	// These variables take precedence over the same variable declared
	// within a variable file.
	TemplateVars map[string]string `hcl:"template_variables,optional"`

	// The nomad job template file to render the template with. If omitted,
	// the plugin will look for a single `*.nomad` file in the current
	// working directory.
	TemplateFile string `hcl:"template_file,optional"`

	// The variable files to render the template with.
	VariableFiles []string `hcl:"variable_files,optional"`

	// This option makes Levant load the Vault token from the current ENV.
	Vault bool `hcl:"vault,optional"`
}

type Platform struct {
	config DeployConfig
}

// Implement Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// Implement Builder
func (p *Platform) DeployFunc() interface{} {
	// return a function which will be called by Waypoint
	return p.deploy
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet

// In addition to default input parameters the registry.Artifact from the Build step
// can also be injected.
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Platform) deploy(
	ctx context.Context,
	src *component.Source,
	job *component.JobInfo,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	var err error
	var result Deployment

	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	result.Id = id
	result.Name = strings.ToLower(fmt.Sprintf("%s-%s", src.App, id))

	if b.config.TemplateVars == nil {
		b.config.TemplateVars = map[string]string{}
	}

	// Add waypoint deploy and artifact template variables so they can be
	// used in the nomad job template file
	b.config.TemplateVars["Workspace"] = job.Workspace
	b.config.TemplateVars["InputDockerImageFull"] = img.Name()
	b.config.TemplateVars["InputDockerImageName"] = img.Image
	b.config.TemplateVars["InputDockerImageTag"] = img.Tag
	b.config.TemplateVars["DeploymentId"] = result.Id
	b.config.TemplateVars["DeploymentApp"] = src.App
	b.config.TemplateVars["DeploymentName"] = result.Name

	u := ui.Status()
	defer u.Close()
	u.Update("Deploy application")

	config := &levant.DeployConfig{
		Client:   &structs.ClientConfig{},
		Deploy:   &structs.DeployConfig{},
		Plan:     &structs.PlanConfig{},
		Template: &structs.TemplateConfig{},
	}

	config.Client.Addr = b.config.Address
	config.Client.ConsulAddr = b.config.ConsulAddress
	config.Deploy.EnvVault = b.config.Vault
	config.Template.VariableFiles = b.config.VariableFiles

	if b.config.TemplateFile != "" {
		config.Template.TemplateFile = b.config.TemplateFile
	} else {
		if config.Template.TemplateFile = helper.GetDefaultTmplFile(); config.Template.TemplateFile == "" {
			err = fmt.Errorf("template_file missing and no default template found.")
			return nil, err
		}
	}

	config.Template.Job, err = template.RenderJob(config.Template.TemplateFile,
		config.Template.VariableFiles, config.Client.ConsulAddr, &b.config.TemplateVars)
	if err != nil {
		u.Step(terminal.StatusError, fmt.Sprintf("[ERROR] levant/command: %v", err))
		return nil, err
	}

	for _, taskGroup := range config.Template.Job.TaskGroups {
		for _, task := range taskGroup.Tasks {
			if task.Env == nil {
				task.Env = map[string]string{}
			}
			for k, v := range b.config.StaticEnvVars {
				task.Env[k] = v
			}
			if taskGroup.Networks != nil {
				for k, v := range deployConfig.Env() {
					task.Env[k] = v
				}
			}
		}
	}

	// Set our ID on the meta.
	config.Template.Job.SetMeta(metaId, result.Id)
	config.Template.Job.SetMeta(metaNonce, time.Now().UTC().Format(time.RFC3339Nano))

	success := levant.TriggerDeployment(config, nil)
	if !success {
		err = fmt.Errorf("Unable to complete deployment.")
		return nil, err
	}
	u.Step(terminal.StatusOK, "Deployment successfully rolled out!")

	// Make sure result.Name matches the job name so that the job can be
	// destroyed with Waypoint even if the levant template uses a different
	// name than the one we provide
	result.Name = *config.Template.Job.Name

	return &result, nil
}
