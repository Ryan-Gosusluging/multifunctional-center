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
    cfg *config.Settings
    generator *client.Generator
}

func New(cfg *config.Settings) *MFC {
    return &MFC{
        cfg: cfg,
        generator: client.NewGenerator(),
    }
}

func (m *MFC) Work(ctx context.Context) {
    clientChan := make(chan client.Client)
    windowChan := make(chan int, m.cfg.WindowCount) // Канал с номерами окон
    for i := 1; i <= m.cfg.WindowCount; i++ {
        windowChan <- i
    }
    // Инициализация семафора
    sem := make(chan struct{}, m.cfg.WindowCount)
    for i := 0; i < m.cfg.WindowCount; i++ {
        sem <- struct{}{}
    }
    var wg sync.WaitGroup
    // Генератор клиентов (горутина)
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
    wg.Add(1)
    go func() {
        defer wg.Done()
        for c := range clientChan {
            <-sem // Захватываем семафор (свободное окно)
            windowNum := <-windowChan // Получаем номер окна
            wg.Add(1)
            go func(client client.Client, wn int) {
                defer func() {
                    sem <- struct{}{} // Освобождаем семафор
                    windowChan <- wn  // Возвращаем номер окна
                    wg.Done()
                }()

                serviceTime := utils.RandomDuration(m.cfg.MinServiceTime, m.cfg.MaxServiceTime)
                fmt.Printf("Окно %d: Клиент %d начал обслуживание (%v)\n", 
                    wn, client.Number, serviceTime)

                select {
                case <-time.After(serviceTime):
                    fmt.Printf("Окно %d: Клиент %d обслужен\n", wn, client.Number)
                case <-ctx.Done():
                    fmt.Printf("Окно %d: Обслуживание клиента %d прервано\n", wn, client.Number)
                }
            }(c, windowNum)
        }
    }()
    wg.Wait()
    close(windowChan)
    close(sem)
    fmt.Println("МФЦ закрылся")
}
