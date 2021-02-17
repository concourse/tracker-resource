ARG base_image=alpine:latest
ARG builder_image=concourse/golang-builder

FROM ${builder_image} as builder
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /assets/out github.com/concourse/tracker-resource/out/cmd/out
RUN go build -o /assets/in github.com/concourse/tracker-resource/in/cmd/in
RUN go build -o /assets/check github.com/concourse/tracker-resource/check/cmd/check
RUN set -e; mkdir /tests; for pkg in $(go list ./...); do \
                cp -a $(go list -f '{{.Dir}}' $pkg) /tests/$(basename $pkg); \
		go test -o "/tests/$(basename $pkg)/run" -c $pkg; \
	done

FROM ${base_image} AS resource
RUN apk add --update bash tzdata ca-certificates git
COPY --from=builder /assets /opt/resource

FROM resource AS tests
COPY --from=builder /tests /tests
RUN set -e; for test in /tests/*/run; do \
                cd $(dirname $test); \
		./run; \
	done

FROM resource
