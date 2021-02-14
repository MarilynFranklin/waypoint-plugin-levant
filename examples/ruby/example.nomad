job "[[ .DeploymentName ]]" {
  datacenters = ["dc1"]
  type = "service"

  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "3m"
    progress_deadline = "10m"
    auto_revert = false
    canary = 0
  }

  migrate {
    max_parallel = 1
    health_check = "checks"
    min_healthy_time = "10s"
    healthy_deadline = "5m"
  }

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

    restart {
      attempts = 2
      interval = "30m"
      delay = "15s"
      mode = "fail"
    }

    ephemeral_disk {
      size = 300
    }

    task "[[ .DeploymentName ]]" {
      driver = "docker"

      config {
        image = "[[ .InputDockerImageFull ]]"

        ports = ["waypoint"]
      }

      resources {
        cpu    = "[[ .cpu ]]"
        memory = "[[ .memory ]]"
      }
    }
  }

  group "[[ .DeploymentApp ]]-worker" {
    count = 1

    restart {
      attempts = 2
      interval = "30m"
      delay = "15s"
      mode = "fail"
    }

    ephemeral_disk {
      size = 300
    }

    task "[[ .DeploymentApp ]]-worker-[[ .DeploymentId ]]" {
      driver = "docker"

      config {
        image = "[[ .InputDockerImageFull ]]"
        command = "worker"
      }

      resources {
        cpu    = "[[ .cpu ]]"
        memory = "[[ .memory ]]"
      }
    }
  }
}
