package main

import (
    "log"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
    "github.com/joho/godotenv"
)

//const numQuotes int = 17
const db_file string = "pbot.db"

const create_query = `
    PRAGMA foreign_keys = ON;    

    CREATE TABLE IF NOT EXISTS qsources (
    qsrc TEXT NOT NULL, 
    PRIMARY KEY(qsrc)
    );
    
    CREATE TABLE IF NOT EXISTS quotes (
    id INTEGER NOT NULL,
    qsrc TEXT NOT NULL,
    quote TEXT,
    PRIMARY KEY(id),
    FOREIGN KEY(qsrc) REFERENCES qsources(qsrc)
    );`

type Quote struct {
    Name string
    Quote string
}

func main(){
    err := godotenv.Load(".env");
    if err != nil{
        panic("couldn't load env");
    }
    var token string = os.Getenv("token");

    dg, err := discordgo.New("Bot " + token)
    if err != nil{
        fmt.Println("error creating discord, ", err)
        return
    }

    dg.AddHandler(messageCreate)
    dg.Identify.Intents = discordgo.IntentsGuildMessages

    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection, ", err)
        return
    }

    fmt.Println("Bot is now running...")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    quotesdb, err := setupDB()
    if err != nil{
        panic(err)
    }

    if m.Content[0] == '!'{
        command := strings.Split(m.Content[1:], " ")
        if command[0] == "quote"{
            var output string = quoteCommand(command, quotesdb)
            if len(output) > 0 {
                s.ChannelMessageSend(m.ChannelID, output)
            }
        }
    }

    quotesdb.Close()
}

func quoteCommand(command []string, quotesdb *sql.DB) string {
    if len(command) < 2{
        return "Missing command argument" 
    }

    if command[1] == "add"{
        return addQuote(command, quotesdb)
    } else if command[1] == "create"{
        return createQuote(command, quotesdb)
    } else {
        row := quotesdb.QueryRow("SELECT quote FROM quotes WHERE qsrc=? ORDER BY RANDOM() LIMIT 1;", command[1])
        randomQuote := Quote{}
        var err error
        if err = row.Scan(&randomQuote.Quote); err == sql.ErrNoRows {
            return fmt.Sprintf("Could not find any quotes by %s", command[1]);
        }
        return randomQuote.Quote
    }

}

func addQuote(command []string, quotesdb *sql.DB) string{
    if len(command) < 4 {
        return "Missing command argument";
    }

    var query string = `
        INSERT INTO quotes(qsrc, quote)
        VALUES(?, ?);
    `
    _, err := quotesdb.Exec(query,
                            command[2], strings.Join(command[3:], " "));

    if err != nil {
        log.Println(err)
        return fmt.Sprintf("Failed to add quote to %s :/", command[2]);
    }
    return "Successfully added quote!"
}

func createQuote(command []string, quotesdb *sql.DB) string{
    if len(command) <  3{
        return "Missing command argument";
    }

    _, err := quotesdb.Exec("INSERT INTO qsources VALUES(?);",
                            command[2]);
    if err != nil{
        return fmt.Sprintf("Error entering %s into database.", command[2]);
    }
    return "Successfully added source";
}

func setupDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", db_file)
    if err != nil {
        return nil, err
    }

    if _, err := db.Exec(create_query); err != nil {
        return nil, err
    }
    return db, nil
}


/*
    var bazinga_quotes []string = []string{
    "Good morning everyone and welcome to 'Science and Society'. I'm Dr. Sheldon Cooper, BS, MS, MA, Ph.D., and ScD. OMG, right?",
    "No cuts, no buts, no coconuts.",
    "So, just to clarify, when you say three, do we stand up or do we pee?",
    "Cause of Injury: Lack of Adhesive Ducks.",
    "A neutron walks into a bar and asks how much for a drink. The bartender replies 'For you, no charge'.",
    "I'm sorry, coffee's out of the question. When I moved to California I promised my mother that I wouldn't start doing drugs.",
    "Well, if you want to see less of me, maybe we should go out again.",
    "I ordered it before you had surgery. It's the urn I was going to put you in.",
    "I'm exceedingly smart. I graduated college at fourteen. While my brother was getting an STD, I was getting a Ph.D.",
    "I'm Batman! Ssssh!",
    "Mom smokes in the car. Jesus is okay with it, but we can't tell dad.",
    "Was the starfish wearing boxer shorts? Because you might have been watching Nickelodeon.",
    "That was tricky because when it comes to alcohol, she generally means business.",
    "'Not knowing is part of the fun.' Was that the motto of your community college?",
    "I would have been here sooner but the bus kept stopping for other people to get on it.",
    "One cries because one is sad. For example, I cry because others are stupid, and that makes me sad.",
    "You're afraid of insects and women. Ladybugs must render you catatonic.",
    }

    var query string = `
        INSERT INTO quotes(qsrc, quote)
        VALUES(?, ?);
    `

    for _,q := range bazinga_quotes{
        _, err := db.Exec(query,
                                "bazinga", q);

        if err != nil {
            panic(err)
        }
        
    }
*/
