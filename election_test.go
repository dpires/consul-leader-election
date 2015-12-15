package main

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestLeaderElection(t *testing.T) {

	g := Goblin(t)
	g.Describe("LeaderElection", func() {
		g.It("can become leader", func() {
			le := &LeaderElection{LeaderKey: "service/test/leader"}
			le.ElectLeader()
			g.Assert(le.IsLeader()).IsTrue()
		})
	})

}
