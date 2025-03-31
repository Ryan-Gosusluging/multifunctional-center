package mfc

import (
    "context"
    "fmt"
    "github.com/Ryan-Gosusluging/multifunctional-center/client"
    "github.com/Ryan-Gosusluging/multifunctional-center/config"
    "github.com/Ryan-Gosusluging/multifunctional-center/utils"
    "sync"
    "time"
)

type MFC struct {
    cfg       *config.Settings
    generator *client.Generator
}

func New(cfg *config.Settings) *MFC {
    return &MFC{
        cfg:       cfg,
        generator: client.NewGenerator(),
    }
}

func (m *MFC) Run(ctx context.Context) {
    clientChan := make(chan client.Client)
    sem := make(chan struct{}, m.cfg.WindowCount)
    
    // Инициализация семафора
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
                c := m.generator.Generate()
                fmt.Printf("Клиент %d пришел\n", c.Number)
                clientChan <- c
            }
        }
    }()

    // Обработчик
    wg.Add(1)
    go func() {
        defer wg.Done()
        for c := range clientChan {
            <-sem
            
            wg.Add(1)
            go func(currentClient client.Client) {
                defer wg.Done()
                defer func() { sem <- struct{}{} }()

                serviceTime := utils.RandomDuration(m.cfg.MinServiceTime, m.cfg.MaxServiceTime)
                fmt.Printf("Клиент %d начал обслуживание (%v)\n", currentClient.Number, serviceTime)

                select {
                case <-time.After(serviceTime):
                    fmt.Printf("Клиент %d обслужен\n", currentClient.Number)
                case <-ctx.Done():
                    fmt.Printf("Обслуживание клиента %d прервано\n", currentClient.Number)
                }
            }(c)
        }
    }()

    wg.Wait()
    close(sem)
    fmt.Println("МФЦ закрылся")
}
