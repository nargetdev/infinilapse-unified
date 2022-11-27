#! /bin/bash
docker buildx build \
            . \
            --push \
            --platform linux/arm64,linux/amd64 \
            --tag nargetdev/infinilapse:1.0.0
