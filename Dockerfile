FROM golang:1.21 AS build 

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /app/fuzzy-user-api cmd/server/main.go


FROM golang:1.21

WORKDIR /app

COPY --from=build /app/fuzzy-user-api /app/fuzzy-user-api

EXPOSE 9205

CMD ["/app/fuzzy-user-api"]

