package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/johansja/talks/lyricplayer"
)

type Player struct {
	playing        bool
	ticker         *time.Ticker
	secondsElapsed int64
	repeat         int
}

func NewPlayer() *Player {
	p := &Player{
		playing: true,
		ticker:  time.NewTicker(1 * time.Second),
	}
	return p
}

func (p *Player) Play() {
	p.playing = true
	go func() {
		for p.playing {
			select {
			case <-p.ticker.C:
				p.secondsElapsed++
				log.Print("Second(s) elapsed: ", p.secondsElapsed)
			default:
				if p.secondsElapsed >= 15 {
					p.Replay()
				}
			}
		}
	}()
}

func (p *Player) Stop() {
	p.ticker.Stop()
	p.playing = false
}

func (p *Player) Replay() {
	p.Stop()
	p.secondsElapsed = 0
	p.ticker = time.NewTicker(1 * time.Second)
	p.repeat++
	p.Play()
}

func (p *Player) GetTime(ctx context.Context, req *player.GetTimeRequest) (*player.GetTimeResponse, error) {
	return &player.GetTimeResponse{
		Time: p.secondsElapsed,
	}, nil
}

func (p *Player) SetTime(ctx context.Context, req *player.SetTimeRequest) (*player.SetTimeResponse, error) {
	p.secondsElapsed = req.Time
	return &player.SetTimeResponse{}, nil
}

func (p *Player) GetLyric(req *player.GetLyricRequest, stream player.Player_GetLyricServer) error {
	curRepeat := p.repeat
	for curRepeat == p.repeat {
		res := &player.GetLyricResponse{
			Lyric: fmt.Sprint("This is lyric at second #", p.secondsElapsed),
		}
		log.Println("Sending lyric: ", res)
		if err := stream.Send(res); err != nil { // HL
			log.Print("Error sending lyric: ", err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	// START TRACE OMIT
	go func() {
		http.ListenAndServe(":2016", nil)
	}()
	// END TRACE OMIT

	p := NewPlayer()
	log.Print("Start playing at: ", time.Now())
	p.Play()

	lis, err := net.Listen("tcp", ":2015")
	if err != nil {
		log.Fatal("Could listen at :2015: ", err)
	}

	s := grpc.NewServer()
	player.RegisterPlayerServer(s, p)
	s.Serve(lis)
}
