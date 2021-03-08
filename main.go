package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"os/signal"
	"syscall"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/bwmarrin/discordgo"
)


var (
	Token string
	ChanID string
	client = http.Client{}
)

func init() {

	flag.StringVar(&Token, "t", "", "BOT Token")
	flag.StringVar(&ChanID, "c", "", "Channel ID")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	
	if m.Content == "!help" {
		msgs := "> [?] List Command :\n`!help`\n`!nhentai <id>`"
		s.ChannelMessageSend( ChanID, msgs)
	}

	if strings.Contains(m.Content, "!nhentai") == true {
		
		cmd := strings.Split(m.Content, " ")
		
		if len(cmd) > 1 {
		
			resp, _ := client.Get("https://nhentai.net/g/" + cmd[1])
			
			defer resp.Body.Close()
			
			data, _ := ioutil.ReadAll(resp.Body)
			
			title := regexp.MustCompile(`(?i)<meta property="og:title" content="(.*?)"`).FindAllSubmatch([]byte(data), -1)
			pic := regexp.MustCompile(`(?i)data-src="(.*?)"`).FindAllSubmatch([]byte(data), -1)
			
			if len(pic) == 0 {
				s.ChannelMessageSend(m.ChannelID, "> [!] ID Not Found")
			}
			
			c := len(pic)
			i := 1
			
			for _, o := range pic {
				reqImg, _ := client.Get(string(o[1]))
				msgs := fmt.Sprintf(">>> ID: %s\nTitle: %s\nPage: %d/%d", cmd[1], title[0][1], i, c)
				s.ChannelFileSendWithMessage( ChanID, msgs, "x.jpg", reqImg.Body)
				i++
			}
			
		} else {
			s.ChannelMessageSend(m.ChannelID, "> [!] Usage: \n`!nhentai 350388`")
		}
	    
	}
	
}
