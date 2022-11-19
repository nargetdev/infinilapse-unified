# Compile stage
FROM golang AS build-env
ENV CGO_ENABLED 0

ARG src_dir=/goSrcFromHost

COPY . $src_dir
WORKDIR $src_dir

#COPY pkg $src_dir
## Copy the predefined netrc file into the location that git depends on
#COPY ./netrc /root/.netrc
#RUN chmod 600 /root/.netrc

RUN go build -gcflags "all=-N -l" -o /webCamEnum $src_dir/pkg/webcamMgmt/cmd
#RUN go build -gcflags "all=-N -l" -o /reformat_epoch_to_hooman ./cmd/reformat_epoch_time_to_hooman/
#RUN go build -gcflags "all=-N -l" -o /compilePriorChunks ./cmd/stitchPrior/

#CMD ["/webcam_cap"]

# Final stage
FROM balenalib/raspberrypi4-64-alpine:latest

# for timezone info make sure to set TZ
RUN apk add --no-cache tzdata

COPY .config/gcloud/application_default_credentials.json /root/.config/gcloud/application_default_credentials.json

# auth for GCP
ENV GOOGLE_APPLICATION_CREDENTIALS /root/.config/gcloud/application_default_credentials.json
RUN apk add ca-certificates

RUN apk add ffmpeg gphoto2 libgphoto2 v4l-utils

COPY --from=build-env /webCamEnum /

# Run
CMD ["/webCamEnum"]