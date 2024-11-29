# Example for streaming speech to text

For this example, we'll use the streaming `SpeechToText` grpc service to transcribe an audio file with `VNG Cloud Speech Text AI`.

## Requirements

- Golang 1.22 or above
- A VNG Cloud Account
- Pakage `pkg-config` and `portaudio` 
## Setup

1. [Enable the Speech-to-Text API](https://ai-speech-text.console.vngcloud.vn/overview)
2. [Create a service account](https://iam.console.vngcloud.vn/service-accounts) with action `ai-stt-tts:StreamingSpeechToText` permission and download the service account key

## Installation

```bash
git clone https://github.com/vngcloud/stt-streaming-example.git
```
The proto file is located at `proto/stt.proto`. You can generate the gRPC client by running the following command:

```bash
cd proto
./generate.sh
cd ../
```

For build and run the example
```bash
#file streaming example
go build -o file ./examples/file/main.go

#mic streaming example
go build -o mic ./examples/mic/main.go
```

## Usage

# file streaming example
```bash
./file -h 
Usage of ./file:
  -client_id string
        Client ID
  -client_secret string
        Client Secret
  -file string
        File path
```
Example:
```bash
./file -client_id <client_id> -client_secret <client_secret> -file <file_path>
```
# mic streaming example
```bash
./mic -h
Usage of ./mic:
  -client_id string
        Client ID
  -client_secret string
        Client Secret
```
Example:
```bash
./mic -client_id <client_id> -client_secret <client_secret>
```