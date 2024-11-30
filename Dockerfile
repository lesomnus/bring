# syntax=docker/dockerfile:1
FROM golang:1.23-alpine3.20 AS builder

ARG TARGETARCH
ARG REPO_NAME
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
			&& apk add --no-cache \
				curl \
				bash \
			&& ./scripts/release.sh \
		; \
	fi


FROM scratch

COPY --from=builder /app/bring /bring
COPY --from=alpine:3.20 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bring"]
