#! /bin/bash
docker buildx build \
            . \
            --platform linux/arm64 \
            --tag nargetdev/infinilapse.chunk-compiler:x-0.5.3
