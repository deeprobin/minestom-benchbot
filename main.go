package main

import (
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/mattn/go-colorable"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	botCount, err := strconv.Atoi(os.Args[0])
	if err != nil {
		log.Fatal("Cannot convert os.Args[0] to int")
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < botCount; i++ {
		go runBot(wg, i)
		// Workaround to not overload the network interface
		if i % 60 == 0 {
			time.Sleep(time.Millisecond * 1000)
		}
	}

	wg.Wait()
	log.Println("Finished.")
}

func runBot(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	log.SetOutput(colorable.NewColorableStdout())
	var c = bot.NewClient()
	c.Auth.Name = "Player" + strconv.FormatInt(int64(i), 10)
	var p = basic.NewPlayer(c, basic.DefaultSettings)

	//Register event handlers
	basic.EventsListener{
		GameStart:  func () error {
			log.Println("Login success: " + c.Name)
			return nil
		},
		Death: func() error {
			log.Printf("Death: " + c.Name)
			_ = p.Respawn()
			return nil
		},
	}.Attach(c)

	// Try to login
	err := c.JoinServer("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	// Join the game
	err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}