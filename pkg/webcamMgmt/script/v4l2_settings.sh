#!/bin/bash

cam=4
amp=1500
v4l2-ctl -d$cam \
-c exposure_absolute=1500 \
-c white_balance_temperature=5000 \
-c exposure_auto=1 \
-c exposure_auto_priority=0 \
-c gain=50 \
--set-fmt-video=width=1920,height=1080 \
--stream-mmap \
--stream-count=1 \
--stream-to=$cam.0-gain.$amp-amplitude.jpg


v4l2-ctl -d24 \
--set-fmt-video=width=1920,height=1080 \
--stream-mmap \
--stream-count=1 \
--stream-to=test.jpg