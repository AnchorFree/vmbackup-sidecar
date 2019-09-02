## binarybuild
##
FROM golang:1.12.9-alpine3.10 as binarybuilder

ENV GO111MODULE on
ENV PROJECT_NAME vmbackup-sidecar

RUN apk --no-cache add git

WORKDIR /go/src/${PROJECT_NAME}
COPY go.mod go.sum /go/src/${PROJECT_NAME}/
RUN go mod download
COPY cmd /go/src/${PROJECT_NAME}/cmd
COPY pkg /go/src/${PROJECT_NAME}/pkg
COPY internal /go/src/${PROJECT_NAME}/internal
RUN cd cmd && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/${PROJECT_NAME}

## awscli + app
##
FROM python:3.7-alpine3.10
RUN pip install --no-cache-dir awscli==1.16.229

# vmbackup app
ENV BINARY vmbackup-sidecar
EXPOSE 8488

WORKDIR /app
COPY --from=binarybuilder /build/${BINARY} bin/${BINARY}

ENTRYPOINT ["/app/bin/vmbackup-sidecar"]