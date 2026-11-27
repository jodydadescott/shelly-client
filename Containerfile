FROM docker.io/library/golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/shelly-client .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/shelly-client /usr/local/bin/shelly-client
ENTRYPOINT ["/usr/local/bin/shelly-client"]
