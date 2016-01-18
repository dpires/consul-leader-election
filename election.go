package election

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulInterface interface {
	GetAgentName() string
	GetKey(string) (*api.KVPair, error)
	PutKey(*api.KVPair) error
	ReleaseKey(*api.KVPair) (bool, error)
	GetSession(string) string
	AquireSessionKey(string, string) (bool, error)
	GetHealthChecks(state string, options *api.QueryOptions) ([]*api.HealthCheck, error)
}

type LeaderElection struct {
	LeaderKey     string
	WatchWaitTime int
	StopElection  chan bool
	Client        ConsulInterface
}

func (le *LeaderElection) CancelElection() {
	le.StopElection <- true
}

func (le *LeaderElection) StepDown() error {
	if le.IsLeader() {
		client := le.Client
		name := client.GetAgentName()
		session := le.GetSession(le.LeaderKey)
		key := &api.KVPair{Key: le.LeaderKey, Value: []byte(name), Session: session}
		released, err := client.ReleaseKey(key)
		if !released || err != nil {
			return err
		} else {
			log.Info("Released leadership")
		}
	}
	return nil
}

func (le *LeaderElection) IsLeader() bool {
	client := le.Client
	name := client.GetAgentName()
	session := le.GetSession(le.LeaderKey)
	kv, err := client.GetKey(le.LeaderKey)
	if err != nil || kv == nil {
		if err != nil {
			log.Error(err)
		}
		log.Info("Leadership key is missing")
		return false
	}

	return name == string(kv.Value) && session == kv.Session
}

func (le *LeaderElection) GetSession(sessionName string) string {
	client := le.Client
	session := client.GetSession(sessionName)
	return session
}

func (le *LeaderElection) ElectLeader() {
	client := le.Client
	name := client.GetAgentName()
	stop := false
	for !stop {
		select {
		case <-le.StopElection:
			stop = true
			log.Info("Stopping election")
		default:
			if !le.IsLeader() {

				session := le.GetSession(le.LeaderKey)

				aquired, err := client.AquireSessionKey(le.LeaderKey, session)

				if aquired {
					log.Infof("%s is now the leader", name)
				}

				if err != nil {
					log.Warn(err)
				}

			}

			kv, err := client.GetKey(le.LeaderKey)

			if err != nil {
				log.Error(err)
			} else {

				if kv != nil && kv.Session != "" {
					log.Info("Current leader: ", string(kv.Value))
					log.Info("Leader Session: ", string(kv.Session))
				}
			}

			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
		}
	}
}
