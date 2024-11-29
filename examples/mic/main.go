package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"stt-streaming-example/helper"
	"stt-streaming-example/stt"
	"sync"

	"github.com/gordonklaus/portaudio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	Len = 2048
)

func main() {
	clientID := flag.String("client_id", "", "Client ID")
	clientSecret := flag.String("client_secret", "", "Client Secret")
	flag.Parse()
	if *clientID == "" || *clientSecret == "" {
		fmt.Println("Please provide clientID and clientSecret")
		return
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	devide, err := portaudio.DefaultInputDevice()
	helper.CheckError(err)
	fmt.Println("Devide: ", devide.Name)
	defer portaudio.Terminate()
	in := make([]int16, Len)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	helper.CheckError(err)
	defer stream.Close()
	helper.CheckError(stream.Start())
	ch := make(chan []int16)
	go streamingSpeechToText(ch, *clientID, *clientSecret)

	fmt.Println("Recording.  Press Ctrl-C to stop.")
	for {
		helper.CheckError(stream.Read())
		ch <- in
		select {
		case <-sig:
			return
		default:
		}
	}
	helper.CheckError(stream.Stop())
	close(ch)

}

func streamingSpeechToText(channel <-chan []int16, clientID, clientSecret string) {
	token := helper.GetVNGCloudToken(clientID, clientSecret)
	conn, err := grpc.NewClient("ai-speech-text.api.vngcloud.vn",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
		// grpc.WithInsecure(),
	)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
	// md.Append("portal-user-id", "11212")
	ctx := metadata.NewOutgoingContext(context.TODO(), md)
	client, err := stt.NewSpeechToTextClient(conn).StreamingSpeechToText(ctx)
	helper.CheckError(err)
	err = client.Send(&stt.StreamingSpeechToTextRequest{
		Request: &stt.StreamingSpeechToTextRequest_Config{
			Config: &stt.StreamingSpeechToTextConfig{
				AudioEncoding:     stt.StreamingSpeechToTextConfig_WAV,
				SampleRateHertz:   16000,
				BytesPerSample:    2,
				AudioChannelCount: 1,
			},
		}})
	helper.CheckError(err)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {

		for data := range channel {
			var buf bytes.Buffer

			err := binary.Write(&buf, binary.LittleEndian, data)
			helper.CheckError(err)
			// fmt.Println(buf.Bytes())
			var chunkType stt.StreamingSpeechToTextData_ChunkType
			if len(data) == Len {
				chunkType = stt.StreamingSpeechToTextData_MIDDLE
			} else {
				chunkType = stt.StreamingSpeechToTextData_LAST
			}

			err = client.Send(&stt.StreamingSpeechToTextRequest{
				Request: &stt.StreamingSpeechToTextRequest_Data{
					Data: &stt.StreamingSpeechToTextData{
						Data:      buf.Bytes(),
						ChunkType: chunkType,
					},
				},
			})
			helper.CheckError(err)
		}
	}()
	for {
		resp, err := client.Recv()
		if err == io.EOF {
			wg.Done()
			break
		}
		helper.CheckError(err)
		log.Println(resp)
	}
	wg.Wait()
	client.CloseSend()

}
