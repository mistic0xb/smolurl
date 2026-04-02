FROM golang:1.26-trixie AS build

WORKDIR /app

COPY go.mod go.sum ./ 

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o /bin/smolurl-server ./cmd/smolurl/main.go

## STAGE 2
FROM scratch 

COPY --from=build /bin/smolurl-server /bin/smolurl-server
COPY --from=build /app/config.yml /app/config.yml

EXPOSE 8080

CMD [ "/bin/smolurl-server" ]