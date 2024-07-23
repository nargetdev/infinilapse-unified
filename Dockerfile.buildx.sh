#! /bin/bash
docker buildx build \
            . \
            --push \
            --platform linux/arm64 \
            --tag nargetdev/infinilapse:arm64-1.1.3
