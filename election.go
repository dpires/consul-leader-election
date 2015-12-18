package election

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulInterface interface {
	GetAgentName() string
	GetKey(string) *api.KVPair
	GetSession(string, *LeaderElection)
	AquireKey(string, string) (bool, error)
}

type LeaderElection struct {
	Session       string
	LeaderKey     string
	WatchWaitTime int
	StopElection  chan bool
	Client       ConsulInterface 
}

func (le *LeaderElection) CancelElection() {
	le.StopElection <- true
}
type GetSessionFunc func(string, *LeaderElection) 

func myfu(name string, le *LeaderElection) {

	client := le.Client
	client.GetSession(name, le)
}

func (le *LeaderElection) IsLeader() bool {
	client := le.Client
	name := client.GetAgentName()
//	le.GetSession(le.LeaderKey)
        le.GetSession(myfu, le.LeaderKey)
	kv := client.GetKey(le.LeaderKey)
	if kv == nil {
		log.Info("Leadership key is missing")
		return false
	}

	return name == string(kv.Value) && le.Session == kv.Session
}
func (le *LeaderElection) GetSession(myFunc GetSessionFunc, sessionName string) {
        
        myFunc(sessionName, le)
}
/*
func (le *LeaderElection) GetSession(sessionName string) {
	client := le.Client
	client.GetSession(sessionName, le)
}
*/

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

				le.GetSession(myfu, le.LeaderKey)//le.LeaderKey)

				aquired, err := client.AquireKey(le.LeaderKey, le.Session)

				if aquired {
					log.Infof("%s is now the leader", name)
				}

				if err != nil {
					panic(err)
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
