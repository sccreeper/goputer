#!/bin/sh

docker build . -t build_goputer

echo """
----------------------------
FINISHED BUILDING CONTAINER
NOW RUNNING MAGE DEV IN CONTAINER
----------------------------
"""

docker run -v "$(pwd)/build:/usr/app/build" -v "$(go env GOCACHE):/root/.cache/go-build" -v "$(pwd)/.parcel-cache:/usr/app/.parcel-cache" build_goputer