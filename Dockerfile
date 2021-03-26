FROM golang:1.14.3-alpine AS build
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download 
RUN GO_ENABLED=0 GOOS=linux go build -o /bin/jwtauthv2 ./cmd/main.go

FROM alpine:latest AS prod

RUN apk add --no-cache bash
COPY --from=build /bin/jwtauthv2 /bin/
EXPOSE 8000
CMD ["/bin/jwtauthv2"]
