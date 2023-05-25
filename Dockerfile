FROM golang AS build

ARG GOOS
ARG GOARCH
ARG GOSUMDB
ARG GOPROXY

ENV CGO_ENABLED=0
ENV GOSUMDB=off
ENV GOPROXY=direct
RUN GOARCH="$(uname -m)"
RUN GOOS="$(uname -s)"
WORKDIR /work
COPY . /work

#RUN go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest
#RUN GOOS=$GOOS GOARCH=$GOARCH go mod download

# Build admission-webhook
RUN GOOS=$GOOS GOARCH=$GOARCH go build -o bin/admission-webhook .

# ---
FROM scratch AS run

COPY --from=build /work/bin/admission-webhook /usr/local/bin/

CMD ["admission-webhook"]
