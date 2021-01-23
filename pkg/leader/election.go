package leader

import (
	"context"
	"flag"
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	leaderElect = flag.Bool("leader-elect", false, "Enable leader election using Consul")
)

type Config struct {
	Key             string
	OnAcquireLeader func()
	ConsulClient    *consulapi.Client
}

type Election struct {
	Config Config

	isLeader  bool
	sessionID string
}

func NewElection(cfg Config) (*Election, error) {
	if cfg.ConsulClient == nil {
		var err error
		cfg.ConsulClient, err = consulapi.NewClient(consulapi.DefaultConfig())
		if err != nil {
			return nil, err
		}
	}

	return &Election{
		Config: cfg,
	}, nil
}

func (e *Election) IsLeader() bool {
	return e.isLeader
}

func (e *Election) Run(ctx context.Context) {
	if !*leaderElect {
		log.Printf("No leader election desired.")
		e.isLeader = true
		go e.Config.OnAcquireLeader()
		return
	}

	log.Printf("Creating Consul session for key %s", e.Config.Key)
	var err error
	e.sessionID, _, err = e.consul().Session().Create(&consulapi.SessionEntry{
		Name: e.Config.Key,
		TTL:  "15s",
	}, nil)
	if err != nil {
		log.Panicf("Error creating Consul session: %v", err)
	}
	go e.consul().Session().RenewPeriodic("15s", e.sessionID, nil, ctx.Done())

	gotLeader, _, err := e.consul().KV().Acquire(&consulapi.KVPair{
		Key:     e.Config.Key,
		Value:   []byte(e.sessionID),
		Session: e.sessionID,
	}, nil)
	if err != nil {
		log.Panicf("Error trying to acquire leadership: %v", err)
	}

	var waitIndex uint64
	for !gotLeader {
		// Watch for the key to change to not have a session anymore, and try to grab leader
		kvPair, meta, err := e.consul().KV().Get(e.Config.Key, &consulapi.QueryOptions{
			WaitIndex: waitIndex,
		})
		if err != nil {
			log.Printf("Error checking on leader key %q: %v", e.Config.Key, err)
			waitIndex = 0
			time.Sleep(10 * time.Second)
			continue
		}

		if meta.LastIndex < waitIndex {
			waitIndex = 0
		} else {
			waitIndex = meta.LastIndex
		}

		if kvPair == nil || kvPair.Session == "" {
			// There's no longer a session, so try to acquire leadership
			gotLeader, _, err = e.consul().KV().Acquire(&consulapi.KVPair{
				Key:     e.Config.Key,
				Value:   []byte(e.sessionID),
				Session: e.sessionID,
			}, nil)
			if err != nil {
				log.Panicf("Error trying to acquire leadership: %v", err)
			}
		}
	}

	// Once we get here, we are the leader
	log.Printf("Started leading")
	e.isLeader = true
	go e.Config.OnAcquireLeader()

	for {
		// Watch to see if our leadership gets revoked
		kvPair, meta, err := e.consul().KV().Get(e.Config.Key, &consulapi.QueryOptions{
			WaitIndex: waitIndex,
		})
		if err != nil {
			log.Panicf("Error checking on leader key %q: %v", e.Config.Key, err)
		}

		if meta.LastIndex < waitIndex {
			waitIndex = 0
		} else {
			waitIndex = meta.LastIndex
		}

		if kvPair.Session != e.sessionID {
			log.Fatalf("Lost leadership")
		}
	}
}

func (e *Election) Stop() {
	if *leaderElect && e.isLeader {
		log.Printf("Releasing leadership voluntarily")
		if _, _, err := e.consul().KV().Release(&consulapi.KVPair{
			Key:     e.Config.Key,
			Value:   []byte{},
			Session: e.sessionID,
		}, nil); err != nil {
			log.Printf("Failed to release lock: %v", err)
		}
	}
}

func (e *Election) consul() *consulapi.Client {
	return e.Config.ConsulClient
}
