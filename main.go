package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hatnote-historical/db"
	"hatnote-historical/handlers"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {

	historyDB, err := db.NewDB("hatnotehistory", "hatnotehistory", "hatnotehistory", "db")
	if err != nil {
		log.Fatalln("Error initialising db")
	}

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
						log.Printf("%s - RCV: %v", lang, c)
						if c.Action == "edit" {
							// Store in db
							sqlStatement := `
              INSERT INTO edits (lang_code, byte_change)
              VALUES ($1, $2);`
							if _, err := historyDB.Exec(sqlStatement, lang, c.ChangeSize); err != nil {
								log.Printf("Error saving edit: %s", err.Error())
							}
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

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
	}))
	h := handlers.New(historyDB, logger)
	app.GET("/edits", h.NetChangePerPeriod)
	app.Run(":8080")

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
