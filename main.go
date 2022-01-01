// ref: https://zenn.dev/saldra/articles/4b4dbca7b8c230

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func loadEnv() {
	// .envの読み込み
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Failed reading .env: %v", err)
	}
	fmt.Println("Loaded .env")
}

func sendMessage(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func sendReply(s *discordgo.Session, channelID string, msg string, reference *discordgo.MessageReference) {
	_, err := s.ChannelMessageSendReply(channelID, msg, reference)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	clientId := os.Getenv("CLIENT_ID")
	user := mc.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", mc.ChannelID, user.Username, user.ID, mc.Content)
	if user.ID != clientId {
		sendMessage(s, mc.ChannelID, user.Mention()+" なんか喋った！")
		sendReply(s, mc.ChannelID, "test", mc.Reference())
	}
}

func main() {
	// 準備
	loadEnv()
	var (
		Token   = "Bot " + os.Getenv("APP_BOT_TOKEN")
		BotName = "<@" + os.Getenv("CLIENT_ID") + ">"
	)
	fmt.Println("Token: ", Token)
	fmt.Println("BotName: ", BotName)

	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Failed to login to Discord.")
		fmt.Println(err)
	}

	// イベントハンドラを追加
	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}
	// 直近の関数(main)の最後に実行される
	defer discord.Close()

	fmt.Println("Listening...")
	// Unixシグナルを受け取るチャネルを作成する
	stopBot := make(chan os.Signal, 1)
	// Unixシグナルを受け取ったらチャネルに通知を送信する
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("aaa")

	// チャネルが値を受信したらここから処理が実行される
	fmt.Println("stopBot received: ", <-stopBot)
	fmt.Println("bbb")
	return
}
