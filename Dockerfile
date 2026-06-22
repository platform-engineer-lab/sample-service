FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /sample-service .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /sample-service /sample-service
EXPOSE 8080
ENTRYPOINT ["/sample-service"]
