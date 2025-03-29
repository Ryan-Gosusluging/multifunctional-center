package mfccore

import (
    "context"
    "fmt"
    "github.com/Ryan-Gosusluging/multifunctional-center/client"
    "github.com/Ryan-Gosusluging/multifunctional-center/config"
    "github.com/Ryan-Gosusluging/multifunctional-center/utils"
    "sync"
    "time"
)

type MFCServer struct {
    cfg      *config.Settings
    clientGen *client.Generator
}

func New(cfg *config.Settings) *MFCServer {
    return &MFCServer{
        cfg:      cfg,
        clientGen: &client.Generator{},
    }
}

func (m *MFCServer) Run(ctx context.Context) {
    clientChan := make(chan *client.Client)
    sem := make(chan struct{}, m.cfg.WindowCount)
    for i := 0; i < m.cfg.WindowCount; i++ {
        sem <- struct{}{}
    }

    var wg sync.WaitGroup

    // Генератор клиентов
    go func() {
        for {
            select {
            case <-ctx.Done():
                close(clientChan)
                return
            default:
                time.Sleep(utils.RandomDuration(m.cfg.MinArrivalTime, m.cfg.MaxArrivalTime))
                client := m.clientGen.Generate()
                fmt.Printf("Клиент %d пришел\n", client.Number)
                clientChan <- client
            }
        }
    }()

    // Обработчик
    wg.Add(1)
    go func() {
        defer wg.Done()
        for client := range clientChan {
            <-sem
            
            wg.Add(1)
            go func(c *client.Client) {
                defer wg.Done()
                defer func() { sem <- struct{}{} }()

                serviceTime := utils.RandomDuration(m.cfg.MinServiceTime, m.cfg.MaxServiceTime)
                fmt.Printf("Клиент %d начал обслуживание (%v)\n", c.Number, serviceTime)

                select {
                case <-time.After(serviceTime):
                    fmt.Printf("Клиент %d обслужен\n", c.Number)
                case <-ctx.Done():
                    fmt.Printf("Обслуживание клиента %d прервано\n", c.Number)
                }
            }(client)
        }
    }()

    wg.Wait()
    close(sem)
    fmt.Println("МФЦ закрылся")
}