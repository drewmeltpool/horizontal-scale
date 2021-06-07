FROM golang:1.16 as build

RUN apt-get update && apt-get install -y ninja-build

WORKDIR /go/src
RUN git clone https://github.com/drewmeltpool/ninja
WORKDIR /go/src/ninja
RUN go get -u ./build/cmd/bood

WORKDIR /go/src/horizontal-scale
COPY . .

RUN CGO_ENABLED=0 bood

# ==== Final image ====
FROM alpine:3.11
WORKDIR /opt/horizontal-scale
COPY entry.sh ./
COPY --from=build /go/src/horizontal-scale/out/bin/* ./
ENTRYPOINT ["/opt/horizontal-scale/entry.sh"]
CMD ["server"]