package discovery

import (
	"fmt"
	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type handler struct {
	joins  chan map[string]string
	leaves chan string
}

func (h *handler) Join(id, addr string) error {
	if h.joins != nil {
		h.joins <- map[string]string{
			"id":   id,
			"addr": addr,
		}
	}
	return nil
}

func (h *handler) Leave(name string) error {
	if h.leaves != nil {
		h.leaves <- name
	}
	return nil
}

func TestMembership(t *testing.T) {
	m, h := setupMember(t, nil)
	m, _ = setupMember(t, m)
	m, _ = setupMember(t, m)

	require.Eventually(t, func() bool {
		return len(h.joins) == 2 && len(m[0].Members()) == 3 && len(h.leaves) == 0
	}, 3*time.Second, 250*time.Millisecond)
	require.NoError(t, m[2].Leave())
	require.Eventually(t, func() bool {
		return len(h.joins) == 2 && len(m[0].Members()) == 3 && m[0].Members()[2].Status == serf.StatusLeft && len(h.leaves) == 1
	}, 3*time.Second, 250*time.Millisecond)
	require.Equal(t, "2", <-h.leaves)
}

func setupMember(t *testing.T, members []*Membership) ([]*Membership, *handler) {
	id := len(members)
	port := 10000 + id
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	tags := map[string]string{
		"rpc_addr": addr,
	}
	cfg := Config{
		NodeName:       fmt.Sprintf("%d", id),
		BindAddr:       addr,
		Tags:           tags,
		StartJoinAddrs: nil,
	}

	h := &handler{}
	if members == nil {
		h.joins = make(chan map[string]string, 3)
		h.leaves = make(chan string, 3)
	} else {
		cfg.StartJoinAddrs = []string{members[0].BindAddr}
	}

	m, err := New(h, cfg)
	require.NoError(t, err)
	members = append(members, m)
	return members, h
}
