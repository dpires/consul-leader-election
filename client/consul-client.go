package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	Client *api.Client
}

func (cc *ConsulClient) GetSession(sessionName string, le *LeaderElection) {
	name := cc.GetAgentName()
	sessions, _, err := cc.Client.Session().List(nil)
	for _, session := range sessions {
		if session.Name == sessionName && session.Node == name {
			le.Session = session.ID
			break
		}
	}
	if le.Session == "" {
		log.Info("No leadership sessions found, creating...")
		sessionEntry := &api.SessionEntry{Name: sessionName}
		le.Session, _, err = cc.Client.Session().Create(sessionEntry, nil)
		if err != nil {
			log.Warn(err)
		}
	}
}

func (cc *ConsulClient) AquireKey(key string, session string, le *LeaderElection) (bool, error) {

	pair := &api.KVPair{
		Key:     key,
		Value:   []byte(cc.GetAgentName()),
		Session: session,
	}

	aquired, _, err := cc.Client.KV().Acquire(pair, nil)

	return aquired, err
}

func (cc *ConsulClient) GetAgentName() string {
	agent, _ := cc.Client.Agent().Self()
	return agent["Config"]["NodeName"].(string)
}

func (cc *ConsulClient) GetKey(keyName string) *api.KVPair {
	kv, _, _ := cc.Client.KV().Get(keyName, nil)
	return kv

}

func (cc *ConsulClient) ReleaseKey(key *api.KVPair) (bool, error) {
	released, _, err := cc.Client.KV().Release(key, nil)
	return release, err
}
