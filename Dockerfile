FROM golang:alpine AS build-env
# RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /src/ /app/
ENTRYPOINT ./goapp
