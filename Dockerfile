ARG GOARCH=arm64
ARG GOOS=linux
FROM golang as builder
WORKDIR /build
COPY ./src ./src
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-w -s" -o ./bin/app ./src/main.go
FROM scratch
WORKDIR /app
COPY --from=builder /build/bin/app /usr/bin/
ENTRYPOINT ["app"]