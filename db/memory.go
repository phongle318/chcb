package db

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/michlabs/fbbot/memory"
)

type persistentMemory struct {
	mutex   *sync.Mutex
	mapping map[string]*persistentStore
}

type persistentStore struct {
	id    string
	mutex *sync.Mutex
	cache map[string]string
}

func newPersistentMemory() *persistentMemory {
	return &persistentMemory{
		mutex:   &sync.Mutex{},
		mapping: make(map[string]*persistentStore),
	}
}

func newPersistentStore() *persistentStore {
	return &persistentStore{
		mutex: &sync.Mutex{},
		cache: make(map[string]string),
	}
}

func (pm *persistentMemory) For(id string) memory.Store {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	ps, ok := pm.mapping[id]
	if !ok {
		ps = newPersistentStore()
		ps.id = id
		pm.mapping[id] = ps

	}
	return ps
}

func (pm *persistentMemory) Delete(id string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	delete(pm.mapping, id)
}

func (ps *persistentStore) Get(key string) string {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	value, ok := ps.cache[key]
	if !ok {
		sender, err := GetSenderByID(ps.id)
		if err == nil {
			ps.cache["customerName"] = sender.FullName
			ps.cache["phone"] = sender.Phone
			ps.cache["gender"] = sender.Gender
		}
		value, ok = ps.cache[key]
	}

	return value
}

func (ps *persistentStore) Set(key string, value string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.cache[key] = value
	sender := Sender{ID: ps.id}
	switch key {
	case "customerName":
		sender.FullName = value
	case "phone":
		sender.Phone = value
	case "gender":
		sender.Gender = value
	default:
		return
	}
	err := UpdateSender(sender)
	if err != nil {
		log.Error("Error when updating sender: ", err)
	}
}

func (ps *persistentStore) Delete(key string) {
	ps.mutex.Lock()

	delete(ps.cache, key)

	switch key {
	case "customerName", "phone", "gender":
		defer ps.Set(key, "")
	}
	defer ps.mutex.Unlock() // prevent deadlock
}

func InitLongTermMemory(bot *fbbot.Bot) {
	bot.LTMemory = newPersistentMemory()
}
