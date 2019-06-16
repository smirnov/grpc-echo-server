package main

import (
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/smirnov/grpc-echo/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	warmupRounds = 200
	rounds       = 1000
)

var (
	address = "localhost:50051"
)

func main() {
	if len(os.Args) > 1 {
		address = os.Args[1]
	}
	// Set up a connection to the server.
	//var durations []time.Duration

	//Warmup run
	min := time.Hour
	max := time.Nanosecond
	sum := time.Nanosecond * 0
	for i := 0; i < warmupRounds; i++ {
		duration, err := establishConnectionAndCallRemote(address)
		if err == nil {
			//append(durations, duration)
			if min > duration {
				min = duration
			}
			if max < duration {
				max = duration
			}
			sum = sum + duration
		}
	}
	avg := sum / warmupRounds
	fmt.Printf("Measurement run for %d warmup rounds\n", warmupRounds)
	fmt.Printf("avg: %s  min: %s  max: %s\n", avg, min, max)

	//Repeat for measurement run
	min = time.Hour
	max = time.Nanosecond
	sum = time.Nanosecond * 0
	for i := 0; i < rounds; i++ {
		duration, err := establishConnectionAndCallRemote(address)
		if err == nil {
			//append(durations, duration)
			if min > duration {
				min = duration
			}
			if max < duration {
				max = duration
			}
			sum = sum + duration
		}
	}
	avg = sum / rounds
	fmt.Printf("Measurement run for %d rounds with new connection for every call\n", rounds)
	fmt.Printf("avg: %s  min: %s  max: %s\n", avg, min, max)

	//Repeat without re-establishing connection
	min = time.Hour
	max = time.Nanosecond
	sum = time.Nanosecond * 0
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)
	for i := 0; i < rounds; i++ {
		start := time.Now()
		_, err = c.Echo(context.Background(), &pb.Message{Msg: "1234"})
		end := time.Now()
		if err == nil {
			duration := end.Sub(start)
			//append(durations, duration)
			if min > duration {
				min = duration
			}
			if max < duration {
				max = duration
			}
			sum = sum + duration
		}
	}
	avg = sum / rounds
	fmt.Printf("Measurement run for %d rounds without re-establishing connection\n", rounds)
	fmt.Printf("avg: %s  min: %s  max: %s\n", avg, min, max)
}

func establishConnectionAndCallRemote(address string) (time.Duration, error) {
	start := time.Now()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n Have you tried appending port ':50051'?", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)
	_, err = c.Echo(context.Background(), &pb.Message{Msg: "1234"})
	end := time.Now()
	return end.Sub(start), err
}
