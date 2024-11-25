# syntax=docker/dockerfile:1
FROM golang:1.23-alpine3.20 AS builder

ARG TARGETARCH
ARG VERSION_NAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bring \
	&& ./bring version

RUN --mount=type=secret,id=github_token,env=GH_TOKEN \
	if [ "$GH_TOKEN" = "" ]; then \
		echo skip push assets to GitHub release \
		; \
	else \
		echo push assets to GitHub release \
		&& eval $(./bring version) \
		&& echo "?" "$BRING_VERSION" "==" "$VERSION_NAME" \
		&& [ "$BRING_VERSION" = "$VERSION_NAME" ] \
		&& apk add --no-cache \
			curl \
		&& curl -sSL "https://github.com/cli/cli/releases/download/v2.62.0/gh_2.62.0_linux_$TARGETARCH.tar.gz" -o ./gh.tar.gz \
		&& tar -xf ./gh.tar.gz \
		&& mv ./gh_*/bin/gh ./gh \
		&& BRING_NAME="./bring-linux-$TARGETARCH" \
		&& cp ./bring "$BRING_NAME" \
		&& ./gh version \
		&& ./gh auth status \
		&& ./gh release --repo lesomnus/bring \
			upload "$VERSION_NAME" "$BRING_NAME" \
			--clobber \
		&& rm -rf \
			./gh.tar.gz \
			./gh_* \
			./bring-* \
		; \
	fi


FROM scratch

COPY --from=builder /app/bring /bring
COPY --from=alpine:3.20 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bring"]
