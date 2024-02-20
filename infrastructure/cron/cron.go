package cron

import (
	"fmt"
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/robfig/cron/v3"
	"health-check/application/interfaces"
	"sync"
	"time"
)

type jobDetails struct {
	id        cron.EntryID
	createdAt time.Time
}

type sCron struct {
	iLogger logger.ILogger

	cron    *cron.Cron
	entries map[uint]jobDetails
	mutex   sync.Mutex
}

func NewCron(logger logger.ILogger) interfaces.ICron {
	return &sCron{
		iLogger: logger,
		cron:    cron.New(),
		entries: make(map[uint]jobDetails),
	}
}

func (r *sCron) AddJob(key uint, createAt time.Time, interval string, job func()) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	entry, ok := r.entries[key]
	if ok {
		if entry.createdAt == createAt {
			return nil
		}
		r.remove(key)
	}

	entryID, err := r.cron.AddFunc(fmt.Sprintf("@every %s", interval), func() {
		defer func() {
			if p := recover(); p != nil {
				r.iLogger.WithAny("panic", p).Error(contextplus.Background(), "recovered from panic")
			}
		}()
		job()
	})
	if err != nil {
		return err
	}

	r.entries[key] = jobDetails{
		id:        entryID,
		createdAt: createAt,
	}

	r.cron.Start()

	return nil
}

func (r *sCron) RemoveJob(key uint) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.remove(key)
}

func (r *sCron) remove(key uint) {
	entry, ok := r.entries[key]
	if ok {
		r.cron.Remove(entry.id)
		delete(r.entries, key)
	}
}
