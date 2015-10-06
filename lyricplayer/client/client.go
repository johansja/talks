package main

import (
	"io"
	"log"
	"time"

	"github.com/johansja/talks/lyricplayer"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func main() {
	// START CONNECT OMIT
	conn, err := grpc.Dial("127.0.0.1:2015", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Couldn't dial: ", err)
	}
	defer conn.Close()
	c := player.NewPlayerClient(conn)
	// END CONNECT OMIT

	ctx := context.Background()

	// START GETTIME OMIT
	log.Print("Getting time")
	gtRes, err := c.GetTime(ctx, &player.GetTimeRequest{}) // HL
	if err != nil {
		log.Fatal("Couldn't get time: ", err)
	}
	log.Print("Second(s) elapsed: ", gtRes.Time)
	// END GETTIME OMIT

	log.Print("Setting time")
	_, err = c.SetTime(ctx, &player.SetTimeRequest{Time: 5})
	if err != nil {
		log.Fatal("Coudln't set time: ", err)
	}
	gtRes, err = c.GetTime(ctx, &player.GetTimeRequest{})
	if err != nil {
		log.Fatal("Couldn't get time: ", err)
	}
	log.Print("Second(s) elapsed: ", gtRes.Time)

	// START GETLYRIC OMIT
	log.Print("Getting 5 lyric")
	glStream, err := c.GetLyric(ctx, &player.GetLyricRequest{}) // HL
	if err != nil {
		log.Fatal("Couldn't get lyric stream: ", err)
	}
	for i := 0; i < 5; i++ {
		glRes, err := glStream.Recv() // HL
		if err != nil {
			log.Fatal("Couldn't get lyric ", i, ": ", err)
		}
		log.Print("Lyric ", i, ": ", glRes.Lyric)
	}
	// END GETLYRIC OMIT

	log.Print("Sleep 2 seconds")
	time.Sleep(2 * time.Second)

	log.Print("Getting until end")
	_, err = c.SetTime(ctx, &player.SetTimeRequest{Time: 5})
	if err != nil {
		log.Fatal("Coudln't set time: ", err)
	}
	for {
		glRes, err := glStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Couldn't get lyric: ", err)
		}
		log.Print("Lyric: ", glRes.Lyric)
	}
}
