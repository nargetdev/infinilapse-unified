#! /bin/bash
docker buildx build \
            . \
            --push \
            --platform linux/arm64 \
            --tag nargetdev/infinilapse:0.8.3
