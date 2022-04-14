/*
    This is the main package for PBot discord bot.

    PBot is a experimental bot for learning discord api as well as golang.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
    "github.com/joho/godotenv"
)

func main(){
    err := godotenv.Load(".env");
    if err != nil{
        panic("couldn't load env");
    }
    var token string = os.Getenv("token");

    sess, err := discordgo.New("Bot " + token)
    if err != nil{
        fmt.Println("error creating discord, ", err)
        return
    }

    pbot := PBot{};
    pbot.setupDB();
    pbot.Session = sess;
    pbot.Session.AddHandler(pbot.messageCreate)
    pbot.Identify.Intents = discordgo.IntentsGuildMessages

    err = pbot.Open()
    if err != nil {
        fmt.Println("error opening connection, ", err)
        return
    }

    fmt.Println("Bot is now running...")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    pbot.Close()
    pbot.qdb.Close();
}
