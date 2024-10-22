FROM golang:1.23-alpine AS build

WORKDIR /todennus-migration

COPY ./todennus-migration/go.mod .
COPY ./todennus-migration/go.sum .

RUN go mod download

COPY . /

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /migrate ./cmd/main.go


FROM scratch

WORKDIR /

COPY --from=build /migrate /
COPY --from=build /todennus-migration/postgres/migration /postgres/migration

ENTRYPOINT ["/migrate", "--env", "", "--path", "/"]
