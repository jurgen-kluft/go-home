package cluster

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/jurgen-kluft/go-home/sensor-server/hollywood/actor"
	"github.com/jurgen-kluft/go-home/sensor-server/hollywood/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Player struct{}

func NewPlayer() actor.Receiver {
	return &Player{}
}

func (p Player) Receive(c *actor.Context) {}

type Inventory struct{}

func NewInventory() actor.Receiver {
	return &Inventory{}
}

func (i Inventory) Receive(c *actor.Context) {}

func TestClusterSelectMemberFunc(t *testing.T) {
	c1, err := New(NewConfig().WithID("A"))
	require.Nil(t, err)
	c2, err := New(NewConfig().WithID("B"))
	require.Nil(t, err)
	c3, err := New(NewConfig().WithID("C"))
	require.Nil(t, err)

	c1.RegisterKind("player", NewPlayer, NewKindConfig())
	c2.RegisterKind("player", NewPlayer, NewKindConfig())
	c3.RegisterKind("player", NewPlayer, NewKindConfig())

	c1.Start()
	c2.Start()
	c3.Start()

	selectMember := func(details ActivationDetails) *Member {
		for _, member := range details.Members {
			if member.ID == "C" {
				return member
			}
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	eventPID := c1.Engine().SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		case ActivationEvent:
			// test that we spawned on member C
			require.Equal(t, c3.Address(), msg.PID.Address)
			cancel()
		case MemberJoinEvent:
			if msg.Member.ID == "C" {
				// Wait till member C is online before activating
				// Activate the actor from member A
				// Which should spawn the actor on member C
				config := NewActivationConfig().WithSelectMemberFunc(selectMember)
				c1.Activate("cancel_receiver", config)
			}
		}
	}, "event")
	c1.Engine().Subscribe(eventPID)
	defer c1.Engine().Unsubscribe(eventPID)

	<-ctx.Done()
	require.Equal(t, context.DeadlineExceeded, ctx.Err())
	c1.Stop()
	c2.Stop()
	c3.Stop()
}

func TestClusterShouldWorkWithDefaultValues(t *testing.T) {
	config := NewConfig()
	c, err := New(config)
	assert.Nil(t, err)
	assert.True(t, len(c.config.id) > 0)
	assert.Equal(t, c.config.region, "default")
}

func TestRegisterKind(t *testing.T) {
	c := makeCluster(t, getRandomLocalhostAddr(), "A", "eu-west")
	c.RegisterKind("player", NewPlayer, NewKindConfig())
	c.RegisterKind("inventory", NewInventory, NewKindConfig())
	assert.True(t, c.HasKindLocal("player"))
	assert.True(t, c.HasKindLocal("inventory"))
}

func TestClusterSpawn(t *testing.T) {
	var (
		c1Addr      = getRandomLocalhostAddr()
		c1          = makeCluster(t, c1Addr, "A", "eu-west")
		c2          = makeCluster(t, getRandomLocalhostAddr(), "B", "eu-west")
		wg          = sync.WaitGroup{}
		expectedPID = actor.NewPID(c1Addr, "player/1")
	)

	wg.Add(2)
	eventPID := c1.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		case MemberJoinEvent:
			if msg.Member.ID == "B" {
				c1.Spawn(NewPlayer, "player", actor.WithID("1"))
			}
		case ActivationEvent:
			assert.True(t, msg.PID.Equals(expectedPID))
			wg.Done()
		}
	}, "event")
	c1.engine.Subscribe(eventPID)

	eventPIDc2 := c2.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		case ActivationEvent:
			assert.True(t, msg.PID.Equals(expectedPID))
			wg.Done()
		}
	}, "event")
	c2.engine.Subscribe(eventPIDc2)

	c1.Start()
	c2.Start()
	wg.Wait()

	c1.Stop()
	c2.Stop()
}

func TestMemberJoin(t *testing.T) {
	c1 := makeCluster(t, getRandomLocalhostAddr(), "A", "eu-west")
	c2 := makeCluster(t, getRandomLocalhostAddr(), "B", "eu-west")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())

	wg := sync.WaitGroup{}
	wg.Add(1)
	eventPID := c1.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		// we do this so we are 100% sure nodes are connected with eachother.
		case MemberJoinEvent:
			fmt.Println(msg)
			if msg.Member.ID == "B" {
				_ = msg
				wg.Done()
			}
		}
	}, "event")
	c1.engine.Subscribe(eventPID)
	c1.Start()
	c2.Start()

	wg.Wait()
	assert.Equal(t, len(c1.Members()), 2)
	assert.True(t, c1.HasKind("player"))

	c1.Stop()
	c2.Stop()
}

func TestActivate(t *testing.T) {
	var (
		addr = getRandomLocalhostAddr()
		c1   = makeCluster(t, addr, "A", "eu-west")
		c2   = makeCluster(t, getRandomLocalhostAddr(), "B", "eu-west")
	)
	c2.RegisterKind("player", NewPlayer, NewKindConfig())

	expectedPID := actor.NewPID(c2.engine.Address(), "player/1")
	wg := sync.WaitGroup{}
	wg.Add(2)
	eventPID := c1.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		// we do this so we are 100% sure nodes are connected with eachother.
		case MemberJoinEvent:
			if msg.Member.ID == "B" {
				// Because c1 doesnt have player registered locally we can only spawned
				// the player on c2
				pid := c1.Activate("player", NewActivationConfig().WithID("1"))
				assert.True(t, pid.Equals(expectedPID))
			}
			wg.Done()
		}
	}, "event")
	c1.engine.Subscribe(eventPID)

	c1.Start()
	c2.Start()

	wg.Wait()
	assert.Equal(t, len(c1.Members()), 2)
	assert.True(t, c1.HasKind("player"))
	assert.True(t, c1.GetActiveByID("player/1").Equals(expectedPID))

	c1.Stop()
	c2.Stop()
}

func TestDeactivate(t *testing.T) {
	addr := getRandomLocalhostAddr()
	c1 := makeCluster(t, addr, "A", "eu-west")
	c2 := makeCluster(t, getRandomLocalhostAddr(), "B", "eu-west")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())

	expectedPID := actor.NewPID(c2.engine.Address(), "player/1")
	wg := sync.WaitGroup{}
	wg.Add(1)
	eventPID := c1.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		case MemberJoinEvent:
			if msg.Member.ID == "B" {
				pid := c1.Activate("player", NewActivationConfig().WithID("1"))
				assert.True(t, pid.Equals(expectedPID))
			}
		case ActivationEvent:
			c1.Deactivate(msg.PID)
		case DeactivationEvent:
			wg.Done()
		}
	}, "event")
	c1.engine.Subscribe(eventPID)

	c1.Start()
	c2.Start()
	wg.Wait()

	assert.Equal(t, len(c1.Members()), 2)
	assert.True(t, c1.HasKind("player"))
	assert.Nil(t, c1.GetActiveByID("player/1"))

	c1.Stop()
	c2.Stop()
}

func TestMemberLeave(t *testing.T) {
	c1Addr := getRandomLocalhostAddr()
	c2Addr := getRandomLocalhostAddr()

	remote := remote.New(c2Addr, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remote))
	if err != nil {
		log.Fatal(err)
	}
	config := NewConfig().
		WithID("B").
		WithRegion("eu-east").
		WithEngine(e)
	c2, err := New(config)
	assert.Nil(t, err)

	c1 := makeCluster(t, c1Addr, "A", "eu-west")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())
	c1.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	eventPID := c1.engine.SpawnFunc(func(c *actor.Context) {
		switch msg := c.Message().(type) {
		case MemberJoinEvent:
			if msg.Member.ID == "B" {
				remote.Stop().Wait()
			}
		case MemberLeaveEvent:
			assert.Equal(t, msg.Member.ID, c2.ID())
			wg.Done()
		}
	}, "event")
	c1.engine.Subscribe(eventPID)
	c2.Start()

	wg.Wait()
	assert.Equal(t, len(c1.Members()), 1)
	assert.False(t, c1.HasKind("player"))

	c1.Stop()
	c2.Stop()
}

func TestMembersExcept(t *testing.T) {
	a := []*Member{
		{
			ID:   "A",
			Host: ":3000",
		},
		{
			ID:   "B",
			Host: ":3001",
		},
	}
	b := []*Member{
		{
			ID:   "A",
			Host: ":3000",
		},
		{
			ID:   "B",
			Host: ":3001",
		},
		{
			ID:   "C",
			Host: ":3002",
		},
	}
	am := NewMemberSet(b...).Except(a)
	assert.Len(t, am, 1)
	assert.Equal(t, am[0].ID, "C")
}

func TestGetActiveByID(t *testing.T) {
	c1Addr := getRandomLocalhostAddr()
	c2Addr := getRandomLocalhostAddr()

	c1 := makeCluster(t, c1Addr, "A", "eu")
	c1.RegisterKind("player", NewPlayer, NewKindConfig())
	c1.Start()

	c2 := makeCluster(t, c2Addr, "B", "eu")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())
	c2.Start()

	pid1 := c1.Activate("player", NewActivationConfig().WithID("1"))
	pid2 := c2.Activate("player", NewActivationConfig().WithID("2"))
	time.Sleep(time.Millisecond * 10)

	pid := c1.GetActiveByID("player/1")
	assert.NotNil(t, pid)
	assert.Equal(t, pid.ID, pid1.ID)

	pid = c1.GetActiveByID("player/2")
	assert.NotNil(t, pid)
	assert.Equal(t, pid.ID, pid2.ID)

	pid = c1.GetActiveByID("player/3")
	assert.Nil(t, pid)

	c1.Stop()
	c2.Stop()
}

func TestGetActiveByKind(t *testing.T) {
	c1Addr := getRandomLocalhostAddr()
	c2Addr := getRandomLocalhostAddr()

	c1 := makeCluster(t, c1Addr, "A", "eu")
	c1.RegisterKind("player", NewPlayer, NewKindConfig())
	c1.Start()

	c2 := makeCluster(t, c2Addr, "B", "eu")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())
	c2.Start()

	pid1 := c1.Activate("player", NewActivationConfig().WithID("1"))
	pid2 := c2.Activate("player", NewActivationConfig().WithID("2"))
	c1.Activate("foo", NewActivationConfig().WithID("2"))
	c1.Activate("bar", NewActivationConfig().WithID("2"))
	time.Sleep(time.Millisecond * 10)

	pids := c1.GetActiveByKind("player")
	assert.Len(t, pids, 2)
	pidsStr := make([]string, 2)
	pidsStr[0] = pids[0].String()
	pidsStr[1] = pids[1].String()
	assert.Contains(t, pidsStr, pid1.String())
	assert.Contains(t, pidsStr, pid2.String())

	c1.Stop()
	c2.Stop()
}

func TestCannotDuplicateActor(t *testing.T) {
	c1Addr := getRandomLocalhostAddr()
	c2Addr := getRandomLocalhostAddr()

	c1 := makeCluster(t, c1Addr, "A", "eu")
	c1.RegisterKind("player", NewPlayer, NewKindConfig())
	c1.Start()

	c2 := makeCluster(t, c2Addr, "B", "eu")
	c2.RegisterKind("player", NewPlayer, NewKindConfig())
	c2.Start()

	pid := c1.Activate("player", NewActivationConfig().WithID("1"))
	time.Sleep(10 * time.Millisecond)
	// Lets make sure we spawn the actor on "our" node. Why?
	// Because when we randomly selected the other node to spawn the actor
	// with the same id on the test will pass.
	// Local registry will prevent duplicated actor IDs from the get go.
	pid2 := c2.Activate("player", NewActivationConfig().WithID("1").WithSelectMemberFunc(func(_ ActivationDetails) *Member {
		return c2.Member()
	}))
	fmt.Println(pid2)
	time.Sleep(time.Millisecond * 10)

	pids := c1.GetActiveByKind("player")
	assert.Len(t, pids, 1)
	assert.Equal(t, pids[0].String(), pid.String())

	c1.Stop()
	c2.Stop()
}

func makeCluster(t *testing.T, addr, id, region string) *Cluster {
	config := NewConfig().
		WithID(id).
		WithListenAddr(addr).
		WithRegion(region)
	c, err := New(config)
	assert.Nil(t, err)
	return c
}

func getRandomLocalhostAddr() string {
	return fmt.Sprintf("127.0.0.1:%d", rand.Intn(50000)+10000)
}
