## binarybuild
##
FROM golang:1.13.0-alpine3.10 as binarybuilder

ENV GO111MODULE on
ENV PROJECT_NAME vmbackup-sidecar

RUN apk --no-cache add git

WORKDIR /go/src/${PROJECT_NAME}
COPY . .
RUN go mod download
RUN VERSION=$(git describe --always --long) && \
    DT=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
    SEMVER=$(git tag --list --sort="v:refname" | tail -n -1) && \
    BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
    cd cmd && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.version=${VERSION} -X main.builddt=${DT} -X main.semver=${SEMVER} -X main.branch=${BRANCH}" -o /build/${PROJECT_NAME}


## awscli + app
##

# awscli has bug syncing empty files under Python3, thus using Python2
# https://github.com/aws/aws-cli/issues/2403
FROM python:2.7-alpine3.10
RUN pip install --no-cache-dir awscli==1.16.236

# Required for full-featured `find` util
RUN apk add findutils

# vmbackup app
ENV BINARY vmbackup-sidecar
EXPOSE 8488

WORKDIR /app
COPY --from=binarybuilder /build/${BINARY} bin/${BINARY}

ENTRYPOINT ["/app/bin/vmbackup-sidecar"]