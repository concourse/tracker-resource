FROM golang:alpine as builder
COPY . /go/src/github.com/concourse/tracker-resource
ENV CGO_ENABLED 0
ENV GOPATH /go/src/github.com/concourse/tracker-resource/Godeps/_workspace:${GOPATH}
ENV PATH /go/src/github.com/concourse/tracker-resource/Godeps/_workspace/bin:${PATH}
RUN go build -o /assets/out github.com/concourse/tracker-resource/out/cmd/out
RUN go build -o /assets/in github.com/concourse/tracker-resource/in/cmd/in
RUN go build -o /assets/check github.com/concourse/tracker-resource/check/cmd/check
RUN set -e; mkdir /tests; for pkg in $(go list ./...); do \
                cp -a $(go list -f '{{.Dir}}' $pkg) /tests/$(basename $pkg); \
		go test -o "/tests/$(basename $pkg)/run" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --update bash tzdata ca-certificates git
COPY --from=builder /assets /opt/resource

FROM resource AS tests
COPY --from=builder /tests /tests
RUN set -e; for test in /tests/*/run; do \
                cd $(dirname $test); \
		./run; \
	done

FROM resource
