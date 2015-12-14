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
	sessions, _, err := client.Session().List(nil)
	sessionID := ""
	for _, session := range sessions {
		if session.Name == leaderKey && session.Node == agent["Config"]["NodeName"] {
			sessionID = session.ID
			break
		}
	}
	if sessionID == "" {
		fmt.Println("No sessions found, getting")
		sessionEntry := &api.SessionEntry{Name: leaderKey}
		sessionID, _, err = client.Session().Create(sessionEntry, nil)
		if err != nil {
			panic(err)
		}
	}

	pair := &api.KVPair{
		Key:     leaderKey,
		Value:   []byte(agent["Config"]["NodeName"].(string)),
		Session: sessionID,
	}
	//	_, err = client.KV().Put(pair, nil)
	// aquire key
	aquired, _, err := client.KV().Acquire(pair, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Aquired:", aquired)
	kv, _, _ := client.KV().Get(leaderKey, nil)
	if kv != nil && kv.Session != "" {
		fmt.Println("Current leader: ", string(kv.Value))
		fmt.Println("Leader Session: ", string(kv.Session))
	}
}
