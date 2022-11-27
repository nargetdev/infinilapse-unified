#! /bin/bash
docker buildx build \
            . \
            --push \
            --platform linux/amd64 \
            --tag nargetdev/infinilapse:amd64-1.0.0
