FROM golang:1.20 AS builder

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build


#FROM registry.access.redhat.com/ubi9-minimal
FROM registry.access.redhat.com/ubi9/ubi

# Create user and group and switch to user's context
RUN dnf -y install shadow-utils \
    && dnf clean all \
    && dnf -y install iputils
RUN useradd --uid 10000 runner
#USER 10000

WORKDIR /app
COPY --from=builder /go/src/app/bin/pinger .

CMD /app/pinger