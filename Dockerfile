FROM golang:1.21-bookworm as build

WORKDIR /app
ADD . /app

ENV APP_NAME stddev-service

RUN go get -d -v ./...
RUN mkdir build
RUN CGO_ENABLED=0 GOOS=linux go build -o build/${APP_NAME} cmd/main.go

# ---------------------------------------------
FROM scratch

ENV APP_NAME stddev-service

COPY --from=build /app/build/${APP_NAME} /

ENTRYPOINT [ "/stddev-service" ]
