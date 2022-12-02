#! /bin/bash
docker buildx build \
            ../.. \
            --file Dockerfile \
            --push \
            --platform linux/arm64,linux/amd64 \
            --tag nargetdev/infinilapse:1.0.4
