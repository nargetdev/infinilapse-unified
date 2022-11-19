#! /bin/bash
docker buildx build \
            -f pkg/webcamMgmt/Dockerfile.webcamMgmt \
            . \
            --push \
            --platform linux/arm64 \
            --tag nargetdev/infinilapse-webcammgmt-list-devices:0.0.0
