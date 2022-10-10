FROM golang:1.19-alpine as dev
WORKDIR /app
COPY . .
RUN mkdir /gocache/
ENV GOCACHE /gocache/
ENV CGO_ENABLED 0
RUN go mod download

FROM dev as build
RUN go build -o monitor cmd/monitor/main.go


FROM alpine as app
COPY --from=build /app/monitor /usr/bin/monitor
CMD '/usr/bin/monitor'