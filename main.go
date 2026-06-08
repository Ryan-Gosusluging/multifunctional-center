package main

import (
    "context"
    "github.com/Ryan-Gosusluging/multifunctional-center/config"
    "github.com/Ryan-Gosusluging/multifunctional-center/mfc"
    "math/rand"
    "time"
	"fmt"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    cfg := config.GetDefault()
    mfc := mfc.New(cfg)
    ctx, cancel := context.WithTimeout(context.Background(), cfg.WorkDuration)
    defer cancel()
    fmt.Println("МФЦ начинает работу")
    mfc.Work(ctx)
}