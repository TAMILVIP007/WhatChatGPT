# Set the base image to Ubuntu 
FROM ubuntu:latest 

# Install packages required for WhatChatGPT application 
RUN apt-get update && \
    apt-get install -y build-essential git wget gcc &&\
    apt-get clean

# Download and install Go 
ENV GO_VERSION 1.19
RUN wget https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz &&\
    tar -xzf go$GO_VERSION.linux-amd64.tar.gz &&\
    chown -R root:root ./go &&\
    mv go /usr/local

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Clone source code from github repository
RUN git clone https://github.com/TAMILVIP007/WhatChatGPT.git

WORKDIR /WhatChatGPT

# Run Go mod tidy and install dependencies
RUN go mod tidy && \ 
    go install github.com/TAMILVIP007/WhatChatGPT

# Compile the Go application
RUN go build -o whatchatgpt main.go

# Make executable file
RUN chmod +x whatchatgpt

# Run the executable
CMD ["./whatchatgpt"]