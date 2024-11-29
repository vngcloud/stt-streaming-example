package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"os"
	"stt-streaming-example/helper"
	"stt-streaming-example/stt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	clientID := flag.String("client_id", "", "Client ID")
	clientSecret := flag.String("client_secret", "", "Client Secret")
	file := flag.String("file", "", "File path")

	flag.Parse()
	if *clientID == "" || *clientSecret == "" {
		log.Panic("Client ID and Client Secret are required")
	}
	if *file == "" {
		log.Panic("File path is required")
	}
	token := helper.GetVNGCloudToken(*clientID, *clientSecret)
	conn, err := grpc.NewClient("ai-speech-text.api.vngcloud.vn",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})),
	)
	helper.CheckError(err)
	defer conn.Close()

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
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
	f, err := os.Open(*file)
	helper.CheckError(err)
	defer f.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				wg.Done()
				return
			}
			helper.CheckError(err)
			var chunkType stt.StreamingSpeechToTextData_ChunkType
			if n == 1024 {
				chunkType = stt.StreamingSpeechToTextData_MIDDLE
			} else {
				chunkType = stt.StreamingSpeechToTextData_LAST
			}

			err = client.Send(&stt.StreamingSpeechToTextRequest{
				Request: &stt.StreamingSpeechToTextRequest_Data{
					Data: &stt.StreamingSpeechToTextData{
						Data:      buf[:n],
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
