#! /bin/bash

IMAGE_TAG=nargetdev/infinilapse:d-1.0.0

echo BUILDING IMAGE: $IMAGE_TAG
docker build \
            . \
            --tag $IMAGE_TAG \

echo PUSHING IMAGE: $IMAGE_TAG

docker push $IMAGE_TAG