package utils

import (
	"math/rand"
	"sync"
	"time"
)

var randMu, randSrc = &sync.Mutex{}, rand.NewSource(time.Now().UnixNano())

func RandomDuration(min, max time.Duration) time.Duration {
	randMu.Lock()
	defer randMu.Unlock()
	if min >= max {
		return min
	}
	r := rand.New(randSrc)
	return min + time.Duration(r.Int63n(int64(max-min)))
}
