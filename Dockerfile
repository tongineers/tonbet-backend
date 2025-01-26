ARG BUILDER_IMAGE="golang:1.23-alpine"
ARG RUNNER_IMAGE="alpine:3.21"

# build stage
FROM ${BUILDER_IMAGE} AS build-stage

ARG GOPRIVATE=""
ARG SSH_KEY

ENV GO111MODULE on
ENV GOPRIVATE ${GOPRIVATE}

RUN apk add --update gcc g++ openssh git make

RUN mkdir -p ~/.ssh/ &&\
    echo "${SSH_KEY}" >> ~/.ssh/id_rsa && chmod 600 ~/.ssh/id_rsa && ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    git config --global url."git@github.com:".insteadOf "https://github.com/"

WORKDIR /app
COPY . .

RUN make build

# production stage
FROM ${RUNNER_IMAGE} AS production-stage

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories &&\
    apk update &&\
    apk add --no-cache\
    ca-certificates

# add pre-compiled TON binaries (fift)
ADD third_party/ton-binaries/fift/fiftlib /usr/local/lib/fiftlib
ADD third_party/ton-binaries/fift/ubuntu-22-0.4.0/fift /usr/local/bin
RUN chmod +x /usr/local/bin/fift
# add static files
ADD scripts/fift scripts/fift

COPY --from=build-stage /app/api api
COPY --from=build-stage /app/build/app .

ENTRYPOINT ["./app"]
# starts app with a "serve" command
CMD ["serve"]