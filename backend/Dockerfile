# prepare builder image
FROM golang:1.23-alpine AS builder


WORKDIR /opt

# copy app source
COPY . .

# install dependencies
RUN go mod download && go mod verify

# build source
RUN go build -tags=viper_bind_struct -o bin/application cmd/main.go

# prepare runner image
FROM alpine AS runner

WORKDIR /opt

# copy executable and dependencies from builder
COPY --from=builder /opt/.env ./
COPY --from=builder /opt/.ssl ./
COPY --from=builder /opt/internal/database/migrations ./internal/database/migrations
COPY --from=builder /opt/bin/application ./
COPY --from=builder /opt/internal/database/mocks/ ./internal/database/mocks

EXPOSE 8080
# run application
CMD ["./application"]
