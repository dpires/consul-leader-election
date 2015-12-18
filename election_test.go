package election

import (
	. "github.com/franela/goblin"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"
)

type FakeClient struct {
	Key string
}

type FakeConsulClient struct {
	Client FakeClient
}

func (fcc *FakeConsulClient) GetAgentName() string {
	return "my node"
}
func (fcc *FakeConsulClient) GetKey(key string) *api.KVPair {
	kv := &api.KVPair{Key: key, Value: []byte("my node"), Session: fcc.Client.Key}
	return kv
}
func (fcc *FakeConsulClient) GetSession(name string, le *LeaderElection) {
	le.Session = name 
}
func (fcc *FakeConsulClient) AquireKey(key string, session string) (bool, error) {
	return true, nil
}

func TestLeaderElection(t *testing.T) {

	g := Goblin(t)
	g.Describe("LeaderElection", func() {
		g.It("can become leader", func() {
			fakeClient := FakeClient{Key: "xservice/leader-election/leader"}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     "xservice/leader-election/leader",
				StopElection:  make(chan bool),
				WatchWaitTime: 1,
				Client:        fake,
			}
			go le.ElectLeader()
			time.Sleep(3 * time.Second)
			le.CancelElection()
			g.Assert(le.IsLeader()).IsTrue()
			//_ = le.StepDown()
		})
		/*
			g.It("can release leadership", func() {
				le := LeaderElection{
					LeaderKey:     "service/leader-election/leader",
					StopElection:  make(chan bool),
					WatchWaitTime: 1,
				}
				go le.ElectLeader()
				time.Sleep(3 * time.Second)
				le.CancelElection()
				g.Assert(le.IsLeader()).IsTrue()
				_ = le.StepDown()
				g.Assert(le.IsLeader()).IsFalse()
			})
		*/
	})

}
