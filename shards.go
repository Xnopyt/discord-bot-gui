package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const identifyRatelimit = time.Second * 5

var shardMan *shardManager
var shardStartWg *sync.WaitGroup

type shardManager struct {
	sync.Mutex

	primarySession *discordgo.Session
	shards         []*discordgo.Session

	shardCount int

	token string

	handlers []interface{}

	running bool
}

func newShardManager(t string) (m *shardManager, err error) {
	m = new(shardManager)
	m.token = t

	m.primarySession, err = discordgo.New(m.token)
	if err != nil {
		return
	}

	resp, err := m.primarySession.GatewayBot()
	if err != nil {
		return nil, err
	}

	m.shardCount = resp.Shards
	if m.shardCount < 1 {
		m.shardCount = 1
	}

	m.running = false

	return
}

func (m *shardManager) addHandler(h interface{}) {
	m.Lock()
	defer m.Unlock()

	m.handlers = append(m.handlers, h)

	if m.running {
		for _, v := range m.shards {
			v.AddHandler(h)
		}
	}
}

func (m *shardManager) start() (err error) {
	m.Lock()
	defer m.Unlock()

	if m.running {
		return errors.New("The shard manager is already running!")
	}

	if m.shardCount < 1 {
		m.shardCount = 1
	}
	log.Println("Attempting to connect to discord using " + strconv.Itoa(m.shardCount) + " shards.")
	shardStartWg.Add(m.shardCount)
	m.shards = make([]*discordgo.Session, m.shardCount)
	for i := range m.shards {
		s, err := discordgo.New(m.token)
		if err != nil {
			return err
		}

		s.ShardCount = m.shardCount
		s.ShardID = i

		for _, v := range m.handlers {
			s.AddHandler(v)
		}

		m.shards[i] = s
	}

	for i := 0; i < m.shardCount; i++ {
		if i != 0 {
			time.Sleep(identifyRatelimit)
		}
		err := m.shards[i].Open()
		if err != nil {
			return err
		}
	}

	m.running = true

	return
}

func (m *shardManager) stop() (err error) {
	m.Lock()
	defer m.Unlock()

	for _, v := range m.shards {
		if x := v.Close(); x != nil {
			err = x
		}
	}

	m.running = false
	return
}

func connectShards(shards int) {
	var err error
	shardMan, err = newShardManager("Bot " + token)
	if err != nil {
		wv.Dispatch(func() {
			wv.Eval("fail()")
			if err.Error() != `HTTP 401 Unauthorized, {"message": "401: Unauthorized", "code": 0}` {
				wv.Eval(`createAlert("` + "Error Creating Sharded Session" + `", '` + err.Error() + `');`)
			}
		})
		return
	}
	shardStartWg = &sync.WaitGroup{}
	shardMan.addHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		log.Println("Shard ID " + strconv.Itoa(s.ShardID) + ": Ready")
		shardStartWg.Done()
	})
	for _, v := range handlers {
		shardMan.addHandler(v)
	}
	if !(shards < 1) {
		shardMan.shardCount = shards
	}
	log.Println("Starting shard manager")
	err = shardMan.start()
	if err != nil {
		shardMan.stop()
		wv.Dispatch(func() {
			wv.Eval("fail()")
			wv.Eval(`createAlert("` + "Error Opening Sharded Session" + `", '` + err.Error() + `');`)
		})
		return
	}
	shardStartWg.Wait()
	ses = shardMan.shards[0]
}
