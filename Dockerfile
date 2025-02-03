# Use the official Golang image with version 1.23
FROM golang:1.23

# Set the working directory
WORKDIR /src

# Install system dependencies
RUN apt update && apt upgrade -y && \
    apt-get install -y qt5-qmake qtbase5-dev libqt5webkit5-dev ffmpeg libasound2-dev

# Copy the Go source code into the container
COPY . .

# Get Go dependencies
RUN go get github.com/aws/aws-sdk-go@v1.55.5 && \
    go get github.com/chromedp/cdproto@v0.0.0-20241003230502-a4a8f7c660df && \
    go get github.com/chromedp/chromedp@v0.10.0 && \
    go get github.com/chromedp/sysutil@v1.0.0 && \
    go get github.com/ebitengine/oto/v3@v3.2.0 && \
    go get github.com/ebitengine/purego@v0.7.1 && \
    go get github.com/gobwas/httphead@v0.1.0 && \
    go get github.com/gobwas/pool@v0.2.1 && \
    go get github.com/gobwas/ws@v1.4.0 && \
    go get github.com/gopherjs/gopherjs@v0.0.0-20190411002643-bd77b112433e && \
    go get github.com/gopxl/beep@v1.4.1 && \
    go get github.com/gopxl/beep/v2@v2.1.0 && \
    go get github.com/jmespath/go-jmespath@v0.4.0 && \
    go get github.com/josharian/intern@v1.0.0 && \
    go get github.com/konsorten/go-windows-terminal-sequences@v1.0.2 && \
    go get github.com/mailru/easyjson@v0.7.7 && \
    go get github.com/pkg/errors@v0.9.1 && \
    go get github.com/sirupsen/logrus@v1.4.1 && \
    go get github.com/therecipe/env_darwin_amd64_513@v0.0.0-20190626001412-d8e92e8db4d0 && \
    go get github.com/therecipe/env_linux_amd64_513@v0.0.0-20190626000307-e137a3934da6 && \
    go get github.com/therecipe/env_windows_amd64_513@v0.0.0-20190626000028-79ec8bd06fb2 && \
    go get github.com/therecipe/env_windows_amd64_513/Tools@v0.0.0-20190626000028-79ec8bd06fb2 && \
    go get github.com/therecipe/qt@v0.0.0-20200904063919-c0c124a5770d && \
    go get github.com/therecipe/qt/internal/binding/files/docs/5.12.0@v0.0.0-20200904063919-c0c124a5770d && \
    go get github.com/therecipe/qt/internal/binding/files/docs/5.13.0@v0.0.0-20200904063919-c0c124a5770d && \
    go get github.com/u2takey/ffmpeg-go@v0.5.0 && \
    go get github.com/u2takey/go-utils@v0.3.1 && \
    go get golang.org/x/crypto@v0.0.0-20200622213623-75b288015ac9 && \
    go get golang.org/x/sys@v0.26.0 && \
    go get golang.org/x/tools@v0.0.0-20190420181800-aa740d480789

# Expose the port your application runs on (if needed)
# EXPOSE 8080 

# Command to run the application
# CMD ["go", "run", "."] # Uncomment and modify if necessary