package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	langCodes := []string{"de", "en"}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logger := log.New(os.Stdout, log.Prefix(), log.Flags())
	// channel used to terminate jobs
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(langCodes))

	for _, langCode := range langCodes {
		go func(lang string) {
			url, err := MakeUrlForLangCode(lang)
			if err != nil {
				log.Printf("Error making URL: %s", err.Error())
				return
			}
			subscriber := &WikiSubscriber{
				url:  url,
				done: done,
				log:  logger,
			}
			results, errors := subscriber.Subscribe()
			go func() {
				defer wg.Done()
				for {
					select {
					case <-done:
						return
					case c := <-results:
						if c.Action == "edit" {
							log.Printf("%s - RCV: %v", lang, c)
						}
					case err := <-errors:
						log.Printf("%s - ERR: %s", lang, err.Error())
					case <-interrupt:
						log.Printf("%s - INTERRUPT", lang)
						close(done)
					}
				}
			}()
		}(langCode)
	}

	wg.Wait()
}

// func main() {
// 	enSubsribe := NewWikiSubscriber("en")
// 	deSubscribe := NewWikiSubscriber("de")
// 	go func() {
// 		if err := enSubsribe.ReadAll(); err != nil {
// 			log.Printf("Error subscribing to wiki edits %s", err.Error())
// 		}
// 		if err := deSubscribe.ReadAll(); err != nil {
// 			log.Printf("Error subscribing to wiki edits %s", err.Error())
// 		}
// 	}()
//
// 	for {
// 		select {
// 		case <-enSubsribe.Done:
// 			log.Println("Done")
// 			return
// 		case r := <-enSubsribe.Results:
// 			log.Printf("RCV: %v", r)
// 		}
// 	}
// }
