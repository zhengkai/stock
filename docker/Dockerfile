# builder stage

FROM golang:latest AS builder

ARG DOCKER_RUNNING=yes
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -yq protobuf-compiler unzip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11

RUN curl -s -L https://github.com/protocolbuffers/protobuf/releases/download/v34.1/protoc-34.1-linux-x86_64.zip -o /usr/local/src/protoc-linux-x86_64.zip \
	&& ls -al /usr/local/src \
	&& unzip -o /usr/local/src/protoc-linux-x86_64.zip -d /usr/local bin/protoc \
	&& unzip -o /usr/local/src/protoc-linux-x86_64.zip -d /usr/local 'include/*'

WORKDIR /project

ARG GOMODCACHE=/cache/mod
ARG GOCACHE=/cache/build

ENV GOMODCACHE=${GOMODCACHE}
ENV GOCACHE=${GOCACHE}

COPY server/go.mod server/go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE} \
	--mount=type=cache,target=${GOCACHE} \
	go mod graph && go mod download && rm -f go.mod go.sum

COPY . .

RUN cd proto && make

ARG GIT_COMMIT
RUN --mount=type=cache,target=${GOMODCACHE} \
	--mount=type=cache,target=${GOCACHE} \
	GIT_COMMIT=${GIT_COMMIT} /project/build/build-server.sh prod

# clean stage

FROM alpine:latest

RUN mkdir /static

# RUN apk add --no-cache asciidoctor cmark tzdata brotli

ENV TZ="Asia/Shanghai"

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /project/dist/prod/stock-server-next /stock-server

# COPY client/dist-git/ /admin

ENV STOCK_WEB=":80"
ENV STOCK_DIR="/static"

EXPOSE 80

CMD ["/stock-server"]
