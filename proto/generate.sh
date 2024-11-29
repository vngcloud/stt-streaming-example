#!/bin/bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
mkdir -p ../stt
protoc --go_out=../stt --go_opt=paths=source_relative \
    --go-grpc_out=../stt --go-grpc_opt=paths=source_relative \
    stt.proto

# change request path /speechtotext.SpeechToText/StreamingSpeechToText -> /speech-api/speechtotext.SpeechToText/StreamingSpeechToText macos
sed -i '' 's/\/speechtotext.SpeechToText\/StreamingSpeechToText/\/speech-api\/speechtotext.SpeechToText\/StreamingSpeechToText/g' ../stt/stt_grpc.pb.go