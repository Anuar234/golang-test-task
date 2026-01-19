FROM alpine:3.19 AS build

WORKDIR /src
RUN apk add --no-cache go ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/server

FROM alpine:3.19
RUN adduser -D app
USER app
COPY --from=build /bin/app /app
EXPOSE 8080
ENV HTTP_ADDR=:8080
ENTRYPOINT ["/app"]
