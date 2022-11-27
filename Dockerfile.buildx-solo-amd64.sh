#! /bin/bash
docker buildx build \
            . \
            --file k3s-manifests/cloud/Dockerfile \
            --push \
            --platform linux/amd64 \
            --tag nargetdev/infinilapse:amd64-1.0.2
