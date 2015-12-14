package main

// Sample leader election implementation from http://consul.io/docs/guides/leader-election.html
import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type ILeaderElect interface {
	GetSession(sessionName string)
	GetConsulClient()
}

type LeaderElect struct {
	Session   string
	LeaderKey string
}

func (le *LeaderElect) GetSession(sessionName string) {
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

func (le *LeaderElect) GetConsulClient() (client *api.Client) {
	config := api.DefaultConfig()
	client, _ = api.NewClient(config)
	return client
}

func main() {
	le := LeaderElect{LeaderKey: "leader-election/leader"}
	client := le.GetConsulClient()
	agent, _ := client.Agent().Self()
	le.GetSession(le.LeaderKey)

	pair := &api.KVPair{
		Key:     le.LeaderKey,
		Value:   []byte(agent["Config"]["NodeName"].(string)),
		Session: le.Session,
	}
	//	_, err = client.KV().Put(pair, nil)
	// aquire key
	aquired, _, err := client.KV().Acquire(pair, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Aquired:", aquired)
	kv, _, _ := client.KV().Get(le.LeaderKey, nil)
	if kv != nil && kv.Session != "" {
		fmt.Println("Current leader: ", string(kv.Value))
		fmt.Println("Leader Session: ", string(kv.Session))
	}
}
