FROM golang:1.19-alpine as dev
WORKDIR /app
COPY . .
RUN mkdir /gocache/
ENV GOCACHE /gocache/
ENV CGO_ENABLED 0
RUN go mod download

FROM dev as build
RUN go build -o indexer cmd/indexer/main.go


FROM alpine as app
COPY --from=build /app/indexer /usr/bin/indexer
CMD '/usr/bin/indexer'