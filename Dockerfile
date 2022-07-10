FROM golang:alpine

EXPOSE 8080
COPY backend /app/backend

ENTRYPOINT ["/app/backend"]

