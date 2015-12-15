package main

// Sample leader election implementation from http://consul.io/docs/guides/leader-election.html
import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type ILeaderElectionion interface {
	GetSession(sessionName string)
	GetConsulClient()
	ElectLeader()
	IsLeader()
}

type LeaderElection struct {
	Session       string
	LeaderKey     string
	WatchWaitTime int
}

func (le *LeaderElection) IsLeader() bool {
	client := le.GetConsulClient()
	agent, _ := client.Agent().Self()
	le.GetSession(le.LeaderKey)
	kv, _, err := client.KV().Get(le.LeaderKey, nil)
	if err != nil {
		fmt.Printf("Unable to check for leadership\n")
		return false
	}
	if kv == nil {
		fmt.Printf("Leadership key is missing")
		return false
	}
	return agent["Config"]["NodeName"] == string(kv.Value) && le.Session == kv.Session
}

func (le *LeaderElection) GetSession(sessionName string) {
	client := le.GetConsulClient()
	agent, _ := client.Agent().Self()
	sessions, _, err := client.Session().List(nil)
	for _, session := range sessions {
		if session.Name == sessionName && session.Node == agent["Config"]["NodeName"] {
			le.Session = session.ID
			break
		}
	}
	if le.Session == "" {
		fmt.Println("No sessions found, getting")
		sessionEntry := &api.SessionEntry{Name: sessionName}
		le.Session, _, err = client.Session().Create(sessionEntry, nil)
		if err != nil {
			panic(err)
		}
	}
}

func (le *LeaderElection) GetConsulClient() (client *api.Client) {
	config := api.DefaultConfig()
	client, _ = api.NewClient(config)
	return client
}

func (le *LeaderElection) ElectLeader() {
	client := le.GetConsulClient()
	agent, _ := client.Agent().Self()

	if !le.IsLeader() {

		le.GetSession(le.LeaderKey)

		pair := &api.KVPair{
			Key:     le.LeaderKey,
			Value:   []byte(agent["Config"]["NodeName"].(string)),
			Session: le.Session,
		}

		aquired, _, err := client.KV().Acquire(pair, nil)

		if aquired {
			fmt.Printf("%s is now the leader\n", agent["Config"]["NodeName"])
		}

		if err != nil {
			panic(err)
		}

	}
	kv, _, _ := client.KV().Get(le.LeaderKey, nil)

	if kv != nil && kv.Session != "" {
		fmt.Println("Current leader: ", string(kv.Value))
		fmt.Println("Leader Session: ", string(kv.Session))
	}

	time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)

	le.ElectLeader()
}

func main() {
	le := LeaderElection{
		LeaderKey:     "service/consul-notifications/leader",
		WatchWaitTime: 10,
	}
	le.ElectLeader()
}
