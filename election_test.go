package election

import (
	"errors"
	. "github.com/franela/goblin"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"
)

type FakeClient struct {
	Key             string
	GetKeyOutput    string
	AquireKeyError  bool
	ReleaseKeyError bool
}

type FakeConsulClient struct {
	Client FakeClient
}

func (fcc *FakeConsulClient) GetHealthChecks(state string, options *api.QueryOptions) ([]*api.HealthCheck, error) {
	return nil, nil
}
func (fcc *FakeConsulClient) GetAgentName() string {
	return "my node"
}

func (fcc *FakeConsulClient) PutKey(key *api.KVPair) error {
	return nil
}

func (fcc *FakeConsulClient) GetKey(key string) (*api.KVPair, error) {
	kv := &api.KVPair{Key: key, Value: []byte("my node"), Session: fcc.Client.Key}
	if fcc.Client.GetKeyOutput == "kv" {
		return kv, nil
	}
	return nil, errors.New("Key not found")
}

func (fcc *FakeConsulClient) ReleaseKey(keyPair *api.KVPair) (bool, error) {
	if fcc.Client.ReleaseKeyError {
		return false, errors.New("ERROR RELEASE KEY")
	}
	fcc.Client.GetKeyOutput = ""
	return true, nil
}

func (fcc *FakeConsulClient) GetSession(name string) string {
	return name
}

func (fcc *FakeConsulClient) AquireSessionKey(key string, session string) (bool, error) {
	if fcc.Client.AquireKeyError {
		return false, errors.New("ERROR")
	}
	fcc.Client.GetKeyOutput = "kv"
	return true, nil
}

func TestLeaderElection(t *testing.T) {

	g := Goblin(t)
	g.Describe("LeaderElection", func() {
		g.It("StepDown() Failure", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader", GetKeyOutput: "kv", ReleaseKeyError: true}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 2,
				Client:        fake,
			}
			go le.ElectLeader()
			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
			le.CancelElection()
			err := le.StepDown()
			isNotNil := (err != nil)
			g.Assert(isNotNil).IsTrue()

		})
		g.It("StepDown()", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader", GetKeyOutput: "kv"}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 2,
				Client:        fake,
			}
			go le.ElectLeader()
			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
			le.CancelElection()
			le.StepDown()
			g.Assert(le.IsLeader()).IsFalse()

		})
		g.It("ElectLeader() Failure", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader", AquireKeyError: true}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 2,
				Client:        fake,
			}
			go le.ElectLeader()
			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
			le.CancelElection()
			g.Assert(le.IsLeader()).IsFalse()
		})
		g.It("ElectLeader()", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader"}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 2,
				Client:        fake,
			}
			go le.ElectLeader()
			time.Sleep(time.Duration(le.WatchWaitTime) * time.Second)
			le.CancelElection()
			g.Assert(le.IsLeader()).IsTrue()
		})
		g.It("CancelElection()", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader", GetKeyOutput: "kv"}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 1,
				Client:        fake,
			}
			go le.ElectLeader()
			le.CancelElection()
			g.Assert(le.IsLeader()).IsTrue()
		})

		g.It("IsLeader()", func() {
			fakeClient := FakeClient{Key: "service/leader-election/leader"}
			fake := &FakeConsulClient{Client: fakeClient}
			le := LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 1,
				Client:        fake,
			}

			g.Assert(le.IsLeader()).IsFalse()

			fakeClient = FakeClient{Key: "service/leader-election/leader", GetKeyOutput: "kv"}
			fake = &FakeConsulClient{Client: fakeClient}
			le = LeaderElection{
				LeaderKey:     fakeClient.Key,
				StopElection:  make(chan bool),
				WatchWaitTime: 1,
				Client:        fake,
			}

			g.Assert(le.IsLeader()).IsTrue()
		})
	})

}
