# build stage
FROM golang:1.25-alpine AS buildstage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/app .

# run stage
FROM alpine:latest
WORKDIR /app

COPY --from=buildstage /bin /bin
COPY --from=buildstage /app/web /app/web

ENV TODO_PASSWORD="12345" \
    TODO_PORT="7540" \
    TODO_DBFILE="scheduler.db"
EXPOSE 7540
CMD ["/bin/app"]