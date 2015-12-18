package election

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulInterface interface {
	GetAgentName() string
	GetKey(string) *api.KVPair
	ReleaseKey(*api.KVPair) (bool, error)
	GetSession(string, *LeaderElection)
	AquireKey(string, string, *LeaderElection) (bool, error)
}

type LeaderElection struct {
	Session       string
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
		le.GetSession(le.LeaderKey)
		key := &api.KVPair{Key: le.LeaderKey, Value: []byte(name), Session: le.Session}
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
	le.GetSession(le.LeaderKey)
	kv := client.GetKey(le.LeaderKey)
	if kv == nil {
		log.Info("Leadership key is missing")
		return false
	}

	return name == string(kv.Value) && le.Session == kv.Session
}
func (le *LeaderElection) GetSession(sessionName string) {
	client := le.Client
	client.GetSession(sessionName, le)
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

				le.GetSession(le.LeaderKey)

				aquired, err := client.AquireKey(le.LeaderKey, le.Session, le)

				if aquired {
					log.Infof("%s is now the leader", name)
				}

				if err != nil {
					log.Warn(err)
				}

			}

			kv := client.GetKey(le.LeaderKey)

			if kv != nil && kv.Session != "" {
				log.Info("Current leader: ", string(kv.Value))
				log.Info("Leader Session: ", string(kv.Session))
			}

			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
		}
	}
}
