package discovery_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
	. "github.com/joshuaejs/godcls/internal/discovery"
	"github.com/stretchr/testify/require"
	dynaport "github.com/travisjeffery/go-dynaport"
)

// TestMembership sets up a cluster with multiple servers and checks that the
// Membership returns all the servers that joined the membership and updates
// after a server leaves the cluster. The handlers join and leave channels
// report how many times each event happened and for what servers. Each member
// has a status:
//   Alive - the server is present and healthy.
//   Leaving - the server is gracefully leaving the cluster.
//   Left - the server gracefully left the cluster.
//   Failed - the server unexpectedly left the cluster.
func TestMembership(t *testing.T) {
	m, handler := setupMember(t, nil)
	m, _ = setupMember(t, m)
	m, _ = setupMember(t, m)
	require.Eventually(t, func() bool {
		return 2 == len(handler.joins) &&
			3 == len(m[0].Members()) &&
			0 == len(handler.leaves)
	}, 3*time.Second, 250*time.Millisecond)
	require.NoError(t, m[2], Leave())
	require.Eventually(t, func() bool {
		return 2 == len(handler.joins) &&
			3 == len(m[0].Members()) &&
			serf.StatusLeft == m[0].Members()[2].Status &&
			1 == len(handler.leaves)
	}, 3*time.Second, 250*time.Millisecond)
	require.Equal(t, fmt.Sprintf("%d", 2), <-handler.leaves)
}

// setupMember is a helper for TestMembership. It sets up a new member under a
// free port and with the member's length as the node name so the names are
// unique. The member's length also indicates whether the member is the
// cluster's initial member, or if there is a cluster to join.
func setupMember(t *testing.T, members []*Membership) (
	[]*Membership, *handler,
) {
	id := len(members)
	ports := dynaport.Get(1)
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
	tags := map[string]string{"rpc_addr": addr}
	c := Config{
		NodeName: fmt.Sprintf("%d", id),
		BindAddr: addr,
		Tags:     tags,
	}
	h := &handler{}
	if len(members) == 0 {
		h.joins = make(chan map[string]string, 3)
		h.leaves = make(chan string, 3)
	} else {
		c.StartJoinAddrs = []string{members[0].BindAddr}
	}
	m, err := New(h, c)
	require.NoError(t, err)
	members = append(members, m)
	return members, h
}

type handler struct {
	joins  chan map[string]string
	leaves chan string
}

// Join
func (h *handler) Join(id, addr string) error {
	if h.joins != nil {
		h.joins <- map[string]string{
			"id":   id,
			"addr": addr,
		}
	}
	return nil
}

// Leave
func (h *handler) Leave(id string) error {
	if h.leaves != nil {
		h.leaves <- id
	}
	return nil
}
