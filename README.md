# Waypoint Plugin Levant

`waypoint-plugin-levant` is an implementation of a platform plugin for
[Waypoint](https://github.com/hashicorp/waypoint) that will allow you to deploy
to a nomad cluster using [levant](https://github.com/hashicorp/levant)
templates.

## Example Usage

After this plugin is installed, you can specify `levant` in the `deploy` clause
of your `waypoint.hcl` file:

```hcl
project = "example-ruby"

app "example-ruby" {
  build {
    use "pack" {}

    registry {
      use "docker" {
        image = "example-ruby"
        tag   = gitrefhash()
        local = true
      }
    }
  }

  deploy {
    use "levant" {}
  }
}
```

By default, the `levant` plugin will look for a single `*.nomad` file in the current
directory to autoload as your nomad job template. It will pass in several
waypoint template variables that you can use when defining your nomad job. (See
the [Template Variables section](#template-variables) for a full list of
available variables):

```hcl
job "[[ .DeploymentName ]]" {
  datacenters = ["dc1"]
  type = "service"

  group "[[ .DeploymentApp ]]" {
    count = 1

    network {
      port "waypoint" {
        to = 5000
      }
    }

    service {
      port = "waypoint"
      name = "[[ .DeploymentApp ]]"
      tags = ["global"]
    }

    task "[[ .DeploymentName ]]" {
      driver = "docker"

      config {
        image = "[[ .InputDockerImageFull ]]"
        ports = ["waypoint"]
      }
    }
  }
}
```

## Quick Start

Take a look at the examples directory to get up and running:

- [Ruby on Rails](examples/ruby/README.md)

## Configuration

These parameters are used in the [`use`
stanza](https://www.waypointproject.io/docs/waypoint-hcl/use) for this plugin.

| Option               | Required   | Type                | Default                 | Description                                                                                                                        |
| -------------------- | ---------- | ------------------- | ----------------------- | ----------------------------------------------------------------------                                                             |
| address              | No         | string              | http://localhost:4646   | The HTTP API endpoint for Nomad where all calls will be made.                                                                      |
| allow_stale          | No         | bool                | false                   | Allow stale consistency mode for requests into nomad.                                                                              |
| canary_auto_promote  | No         | int                 | 0                       | The time period in seconds that Levant should wait for before attempting to promote a canary deployment.                           |
| consul_address       | No         | string              | localhost:8500          | The Consul host and port to use when making Consul KeyValue lookups.                                                               |
| force_batch          | No         | bool                | false                   | Forces a new instance of the periodic job. A new instance will be created even if it violates the job's prohibit_overlap settings. |
| force_count          | No         | bool                | false                   | Use the taskgroup count from the Nomad job file instead of the count that is obtained from the running job count.                  |
| static_environment   | No         | map[string]string   |                         | Environment variables to add to the job.                                                                                           |
| template_variables   | No         | map[string]string   |                         | Variables that are meant to configure the nomad job template file.                                                                 |
| variable_files       | No         | []string            |                         | The variable files to render the template with.                                                                                    |
| vault                | No         | bool                | false                   | This option makes Levant load VAULT_TOKEN from the current ENV.                                                                    |
| prevent_destroy      | No         | bool                | false                   | This option prevents Waypoint from destroying the nomad job.                                                                       |

```hcl
deploy {
  use "levant" {
    template_variables = {
      cpu = 250
      memory = 250
    }

    static_environment = {
      GEM_PATH = "/layers/heroku_ruby/gems/vendor/bundle/ruby/2.6.0"
      PATH = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:bin:/layers/heroku_ruby/gems/vendor/bundle/ruby/2.6.0/bin"
    }
  }
}
```

## Template Variables

Along with any `template_variables` added in the use stanza, this plugin will
also make the following variables available to your nomad job template file:

| Name                 | Description                                                                                                                            |
|----------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| Workspace            | The workspace name that the Waypoint deploy is running in. This lets you potentially deploy to different clusters based on this value. |
| InputDockerImageFull | The full Docker image name and tag.                                                                                                    |
| InputDockerImageName | The Docker image name, without the tag.                                                                                                |
| InputDockerImageTag  | The Docker image tag, such as "latest".                                                                                                |
| DeploymentId         | Generated deployment id                                                                                                                |
| DeploymentApp        | Waypoint App name                                                                                                                      |
| DeploymentName       | Generated deployment name which can be used as the job name                                                                            |
