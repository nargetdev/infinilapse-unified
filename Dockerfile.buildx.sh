#! /bin/bash
docker buildx build \
            . \
            --push \
            --platform linux/arm64 \
            --tag nargetdev/infinilapse:x-0.6.0
