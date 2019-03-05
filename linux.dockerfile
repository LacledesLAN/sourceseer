# escape=`

FROM golang as builder-image

ADD . /usr/local/go/src/github.com/lacledeslan/sourceseer/

WORKDIR /usr/local/go/src/github.com/lacledeslan/sourceseer/

RUN mkdir /output &&`
    go get github.com/jessevdk/go-flags &&`
    go build -o /output/csgotourney ./cmd/csgotourney/

FROM alpine

HEALTHCHECK NONE

ARG BUILDNODE="unspecified"
ARG SOURCE_COMMIT

LABEL maintainer="Laclede's LAN <contact @lacledeslan.com>" `
      com.lacledeslan.build-node=$BUILDNODE `
      org.label-schema.schema-version="1.0" `
      org.label-schema.url="https://github.com/LacledesLAN/sourceseer" `
      org.label-schema.vcs-ref=$SOURCE_COMMIT `
      org.label-schema.vendor="Laclede's LAN" `
      org.label-schema.description="SourceSeer" `
      org.label-schema.vcs-url="https://github.com/LacledesLAN/sourceseer"

COPY --from=builder-image /output /app/

CMD ["--help"]

ENTRYPOINT ["/app/csgotourney"]
