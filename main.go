package main

import (
	"log"
	"fmt"
    "os"
    "strconv"

	"github.com/gorcon/rcon"
    "github.com/bwmarrin/discordgo"
    "strings"
    "slices"
    "time"
)


func parseResponse(text string) []string {
    if strings.HasPrefix(text, "There are 0") {
        return []string{}
    }
    // There are 2 of a max of 20 players online: Steve, Alex
    var playerlist string = strings.Split(text, ":")[1];
    playerlist = strings.TrimSpace(playerlist);

    return strings.Split(playerlist, ", ");
}

func getPlayers(conn *rcon.Conn) []string {
	response, err := conn.Execute("/list")
    if err != nil {
		log.Print(err)
	}
	fmt.Println(response)
    return parseResponse(response)  
}


func sendMessage(dis *discordgo.Session, channel_id string, text string) {
    dis.ChannelMessageSend(channel_id, text)
}

func main() {
    var rcon_host string = os.Getenv("RCON_HOST");

    if rcon_host == "" {
        log.Fatal("Missing RCON_HOST value")
    }

    var rcon_port string = os.Getenv("RCON_PORT");
    if rcon_port == "" {
        log.Fatal("Missing RCON_PORT value")
    }

    _, err := strconv.Atoi(rcon_port)

    if err != nil {
        log.Fatalf("Unable to read RCON_PORT value %s as int", rcon_port)
    }

    var connection_str string = fmt.Sprintf("%s:%s", rcon_host, rcon_port)
    
    var rcon_pass = os.Getenv("RCON_PASS")

	conn, err := rcon.Dial(connection_str, rcon_pass)
	if err != nil {
		log.Fatal(err)
	}

    defer conn.Close()


    var discord_token string = os.Getenv("DISCORD_TOKEN")

    if discord_token == "" {
        log.Fatal("Missing DISCORD_TOKEN value")
    }

    var channel_id string = os.Getenv("DISCORD_CHANNEL")
    
    if channel_id == "" {
        log.Fatal("Missing DISCORD_CHANNEL value")
    }
    discord, err := discordgo.New("Bot " + discord_token)
    
    var oldPlayers []string = getPlayers(conn);
    log.Println("Starting minecraftd...");
    for {
        var newPlayers []string = getPlayers(conn);


        // handle new additions; logon
        for _, player := range newPlayers {
            if !slices.Contains(oldPlayers, player) {
                // Player has logged on
                sendMessage(discord, channel_id, player + " has logged on")
                log.Printf("Player has logged on %s\n", player)
            }
        }

        // handle new reductions; logoff

        for _, player := range oldPlayers {
            if !slices.Contains(newPlayers, player) {
                // Player has logged off

                sendMessage(discord, channel_id, player + " has logged off")
                log.Printf("Player has logged off %s\n", player) 
            }
        }

        oldPlayers = newPlayers
        time.Sleep(5 * time.Second)
    }
}
