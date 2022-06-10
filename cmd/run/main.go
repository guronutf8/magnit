package main

import (
	"magnit/internal/api"
	"magnit/internal/conf"
	"sync"
)

var configuration conf.Config

func main() {
	configuration = conf.Conf()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		// стартуем http сервис
		api.Init(configuration)
		wg.Done()
	}()
	wg.Wait()
}
