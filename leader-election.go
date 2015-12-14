package main

// Sample leader election implementation from http://consul.io/docs/guides/leader-election.html
import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func main() {
	const leaderKey = "leader-election/leader"
	config := api.DefaultConfig()
	client, _ := api.NewClient(config)
	agent, _ := client.Agent().Self()
	// get session
	sessionEntry := &api.SessionEntry{Name: leaderKey}
	sessionID, _, err := client.Session().Create(sessionEntry, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(sessionID)

	pair := &api.KVPair{
		Key:     leaderKey,
		Value:   []byte(agent["Config"]["NodeName"].(string)),
		Session: sessionID,
	}
	//	_, err = client.KV().Put(pair, nil)
	// aquire key
	aquired, _, err := client.KV().Acquire(pair, nil)

	fmt.Println(pp)
	if err != nil {
		panic(err)
	}
	fmt.Println(aquired)
}
