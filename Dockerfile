FROM golang:1.20

WORKDIR /code

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them
# in subsequent builds if they change.
COPY go.mod go.sum* ./
RUN go mod download && go mod verify

CMD ./scripts/run-tests.sh
