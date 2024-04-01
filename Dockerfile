#build stage
# FROM tggolangbot:staging AS builder
FROM shoutertm/dockerhub:tggolangbot AS builder
# RUN apk add --no-cache git
WORKDIR /go/src/tg_golang_bot
COPY . .
# RUN go get -d -v ./...
RUN go build -o /go/bin/tg_golang_bot -v ./cmd/app/main.go

# final stage
FROM alpine:latest
ENV LANGUAGE="en"
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /go/bin/tg_golang_bot /tg_golang_bot
ENTRYPOINT /tg_golang_bot
LABEL Name=tggolangbot Version=0.0.1
EXPOSE 3000
EXPOSE 4000
EXPOSE 80/tcp 