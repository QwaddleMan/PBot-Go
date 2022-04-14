package main

import (
    "log"
	"database/sql"
	"fmt"
    "strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

// The database name.
const db_file string = "pbot.db"

// The database creation query.
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


type PBot struct {
    *discordgo.Session
    qdb *sql.DB
}

// messageParse is the handler function for incomming messages.
func (pbot *PBot) messageCreate(sess *discordgo.Session, mesg *discordgo.MessageCreate) {
    if mesg.Content[0] == '!'{
        command := strings.Split(mesg.Content[1:], " ")
        if command[0] == "quote"{
            var output string = pbot.quoteCommand(command);
            if len(output) > 0 {
                pbot.ChannelMessageSend(mesg.ChannelID, output)
            }
        }
    }
}

// quoteCommand method controls quotes. You can add quotes to sources,
// create new sources and get a random quote from a given source.
func (pbot *PBot) quoteCommand(command []string) string {
    if len(command) < 2{
        return "Missing command argument" 
    }

    var src string = "";
    var quote string = "";

    if len(command) >= 3 {
        src = command[2];
    }

    if len(command) > 3 {
        quote = strings.Join(command[3:], " ");
    }

    if command[1] == "add"{
        if len(command) < 4 {
            return "Missing command argument";
        }
        return pbot.addQuote(src, quote);
    } else if command[1] == "create"{
        if len(command) <  3{
            return "Missing command argument";
        }
        return pbot.createSource(src);
    } else {
        var query string = `
        SELECT qsrc, quote FROM quotes WHERE qsrc=? ORDER BY RANDOM() LIMIT 1;
        `;
        row := pbot.qdb.QueryRow(query, command[1])
        randomQuote := Quote{};
        var err error
        if err = row.Scan(&randomQuote.Quote, &randomQuote.Name); err == sql.ErrNoRows {
            return fmt.Sprintf("Could not find any quotes by %s", command[1]);
        }
        return randomQuote.toString()
    }

}

// addQuote adds a quote to a source.
func (pbot *PBot) addQuote(src string, quote string) string{
    var query string = `
        INSERT INTO quotes(qsrc, quote)
        VALUES(?, ?);
    `
    _, err := pbot.qdb.Exec(query,src, quote);

    if err != nil {
        log.Println(err)
        return fmt.Sprintf("Failed to add quote to %s :/", src);
    }
    return "Successfully added quote!"
}

// createSource creates a source to quote.
func (pbot *PBot) createSource(src string) string{

    _, err := pbot.qdb.Exec("INSERT INTO qsources VALUES(?);", src);
    if err != nil{
        return fmt.Sprintf("Error entering %s into database.", src);
    }
    return "Successfully added source";
}

// setupDB sets up the database and adds it to the pbot structure.
func (pbot *PBot) setupDB() {
    db, err := sql.Open("sqlite3", db_file);
    if err != nil {
        panic(err);
    }

    if _, err := db.Exec(create_query); err != nil {
        panic(err);
    }

    pbot.qdb = db;
}
