FROM golang AS builder

# copy the source
ADD .   /go/src/github.com/Eldarrin/devops-tools/cmd/backend
WORKDIR /go/src/github.com/Eldarrin/devops-tools/cmd/backend

# install dependencies

# build the sample
RUN CGO_ENABLED=0 go build -o /go/bin/backend .

FROM golang:alpine

EXPOSE 8080
COPY --from=builder /go/bin/backend .

ENTRYPOINT ["/app/backend"]

