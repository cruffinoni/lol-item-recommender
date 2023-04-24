package api

import (
	"sync"
	"time"

	"LoLItemRecommender/internal/printer"
)

const (
	maxRatePerSeconds = 20
	maxRatePerMinute  = 100
)

type Bucket struct {
	mx            *sync.RWMutex
	ticker        *time.Ticker
	expiration    time.Duration
	lastT         time.Time
	token         int
	initialTokens int
	done          chan bool
}

func NewBucket(token int, expiration time.Duration) *Bucket {
	b := &Bucket{
		mx:            &sync.RWMutex{},
		lastT:         time.Now(),
		expiration:    expiration,
		token:         token,
		initialTokens: token,
		ticker:        time.NewTicker(expiration),
	}
	go b.resetTokenLoop()
	return b
}

func (b *Bucket) GetToken() int {
	b.mx.RLock()
	defer b.mx.RUnlock()
	return b.token
}

func (b *Bucket) GetWaitingTime() time.Duration {
	b.mx.RLock()
	defer b.mx.RUnlock()

	timeSinceLastConsume := time.Now().Sub(b.lastT)
	timeRemaining := b.expiration - timeSinceLastConsume

	if timeRemaining < 0 {
		return 0
	}
	return timeRemaining
}

func (b *Bucket) TryConsumeToken() bool {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.token-1 >= 0 {
		b.token--
		b.lastT = time.Now()
		return true
	}
	return false
}

func (b *Bucket) CanConsumeToken() bool {
	b.mx.RLock()
	defer b.mx.RUnlock()
	return b.token-1 >= 0
}

func (b *Bucket) resetTokenLoop() {
	for {
		select {
		case <-b.ticker.C:
			b.mx.Lock()
			b.token = b.initialTokens
			b.mx.Unlock()
		case <-b.done:
			return
		}
	}
}

func (b *Bucket) Close() {
	b.ticker.Stop()
	b.done <- true
}

type Rate struct {
	mx              *sync.RWMutex
	totalUsage      int
	perSecondBucket *Bucket
	perMinuteBucket *Bucket
}

func NewRate() *Rate {
	return &Rate{
		mx:              &sync.RWMutex{},
		perSecondBucket: NewBucket(maxRatePerSeconds, 1*time.Second),
		perMinuteBucket: NewBucket(maxRatePerMinute, 1*time.Minute),
	}
}

func (r *Rate) CanConsumeTokens() (bool, time.Duration) {
	r.mx.Lock()
	defer r.mx.Unlock()

	canConsumeSecond := r.perSecondBucket.CanConsumeToken()
	canConsumeMinute := r.perMinuteBucket.CanConsumeToken()
	if canConsumeSecond && canConsumeMinute {
		r.perSecondBucket.TryConsumeToken()
		r.perMinuteBucket.TryConsumeToken()
		r.totalUsage++
		if r.totalUsage%50 == 0 {
			printer.Debug("{-F_RED,BOLD}%d{-RESET} tokens consumed", r.totalUsage)
		}
		return true, 0
	}
	waitTimeSecond := r.perSecondBucket.GetWaitingTime()
	waitTimeMinute := r.perMinuteBucket.GetWaitingTime()
	return false, maxDuration(waitTimeSecond, waitTimeMinute)
}

func (r *Rate) GetTotalUsage() int {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.totalUsage
}
func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
