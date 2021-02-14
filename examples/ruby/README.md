# Waypoint Ruby on Rails Example

A barebones Rails app, which can easily be deployed by Waypoint (from
[waypoint-examples](https://github.com/hashicorp/waypoint-examples/tree/main/docker/ruby)).

## Up and Running

### Install dependencies

- [waypoint](https://learn.hashicorp.com/tutorials/waypoint/get-started-install?in=waypoint/get-started-nomad)

      brew tap hashicorp/tap
      brew install hashicorp/tap/waypoint

- [nomad](https://learn.hashicorp.com/tutorials/nomad/get-started-install?in=nomad/get-started)

      brew tap hashicorp/tap
      brew install hashicorp/tap/nomad

- [waypoint-plugin-levant](https://github.com/MarilynFranklin/waypoint-plugin-levant/releases)

    Navigate to https://github.com/MarilynFranklin/waypoint-plugin-levant/releases.
    After downloading the plugin, unzip the package to this directory.

### Create a nomad environment

    nomad agent -dev -network-interface="en0"


### Install the Waypoint server

    export NOMAD_ADDR='http://localhost:4646'
    waypoint install -platform=nomad -nomad-dc=dc1 -accept-tos


### Initialize Waypoint

    cd examples/ruby
    waypoint init


### Deploy the application with waypoint up

    waypoint up

Waypoint will show the result of your deployment in the Terminal, along with
your specific preview URL.

You can also view the server side of the deployment in the
[Nomad UI](http://localhost:4646/ui/jobs).


### Destroy the deployment

    waypoint destroy
