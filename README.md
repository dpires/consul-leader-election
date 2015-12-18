# consul-leader-election

[![Build Status](https://travis-ci.org/dpires/consul-leader-election.svg?branch=master)](https://travis-ci.org/dpires/consul-leader-election)
[![Coverage Status](https://coveralls.io/repos/dpires/consul-leader-election/badge.svg?branch=master&service=github)](https://coveralls.io/github/dpires/consul-leader-election?branch=master)

An implementation of Consul leader election based on https://www.consul.io/docs/guides/leader-election.html

## Usage

```
import (
	"github.com/dpires/consul-leader-election"
	"github.com/dpires/consul-leader-election/client"
	"github.com/hashicorp/consul/api"
)

config := api.DefaultConfig()                                  // Create a new api client config
consulclient, _ := api.NewClient(config)                       // Create a Consul api client

leaderElection := &election.LeaderElection{
        StopElection:  make(chan bool),                        // The channel for stopping the election
        LeaderKey:     "service/my-leadership-service/leader", // The leadership key to create/aquire
        WatchWaitTime: 10,                                     // Time in seconds to check for leadership
        Client: &client.ConsulClient{Client:consulclient},     // The injected Consul api client
}

go leaderElection.ElectLeader()                                // Run the election
```
