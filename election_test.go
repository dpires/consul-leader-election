package election

import (
	. "github.com/franela/goblin"
	"testing"
	"time"
)

func TestLeaderElection(t *testing.T) {

	g := Goblin(t)
	g.Describe("LeaderElection", func() {
		g.It("can become leader", func() {
			le := LeaderElection{LeaderKey: "service/leader-election/leader", StopElection: make(chan bool), WatchWaitTime: 1}
			go le.ElectLeader()
			time.Sleep(3 * time.Second)
			le.StopElection <- true
			g.Assert(le.IsLeader()).IsTrue()
		})
	})

}
