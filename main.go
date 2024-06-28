package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"layeh.com/gopus"
	"log"
	"os"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	vc, err := dg.ChannelVoiceJoin(os.Getenv("GUILD_ID"), os.Getenv("CHANNEL_ID"), false, false)
	if err != nil {
		fmt.Println("error joining vc,", err)
		return
	}

	data, err := os.ReadFile("bebra.pcm")
	if err != nil {
		fmt.Println("error reading file,", err)
		return
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	buf := bufio.NewReaderSize(bytes.NewReader(data), 16384)

	_ = vc.Speaking(true)

	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(buf, binary.LittleEndian, &audiobuf)

		opus, err := opusEncoder.Encode(audiobuf, frameSize, maxBytes)
		if err != nil {
			return
		}

		vc.OpusSend <- opus
	}

}
