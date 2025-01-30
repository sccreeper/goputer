#!/bin/sh

docker build . -t build_goputer

echo """
----------------------------
FINISHED BUILDING CONTAINER
NOW RUNNING MAGE DEV IN CONTAINER
----------------------------
"""

docker run -v "$(pwd)/build:/usr/app/build" build_goputer