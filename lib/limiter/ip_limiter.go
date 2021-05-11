/*
基于 ip地址的限流器
可以配合middleware使用
*/

package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

// IPLimiter 基于ip的限流
type IPLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPLimiter return a new ip limiter
//	每r秒b个
func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// addIP add ip address to limiter
func (m *IPLimiter) addIP(ip string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	limiter := rate.NewLimiter(m.r, m.b)
	m.ips[ip] = limiter
	return limiter
}

// GetLimiter return
func (m *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	m.mu.Lock()
	limiter, ok := m.ips[ip]
	if !ok {
		m.mu.Unlock()
		return m.addIP(ip)
	}
	m.mu.Unlock()
	return limiter
}
