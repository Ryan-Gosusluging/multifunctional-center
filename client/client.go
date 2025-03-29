package client

import (
    "sync"
)

type Client struct {
    Number int
}

type Generator struct {
    counter int
    mu      sync.Mutex
}

func (g *Generator) Generate() *Client {
    g.mu.Lock()
    defer g.mu.Unlock()
    g.counter++
    return &Client{Number: g.counter}
}