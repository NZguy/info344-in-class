#!/usr/bin/env bash

#TODO:
# - create private network,
# - run hellosvc in network with no published ports
# - run gateway in network publishing port 443
#   and using volumes to give it access to your
#   cert and key files in `gateway/tls`
#   and setting environment variables
#    - CERTPATH = path to cert file in container
#    - KEYPATH = path to private key file in container
#    - HELLOADDR = net address of hellosvc container

# If result is a zero length string
if [ -z "$(docker network ls -q -f name=msdemonet)"]
then
    docker network create msdemonet
fi

docker run -d \
--name hellosvc1 \
--network msdemonet \
duncanandrew/hellosvc

docker run -d \
--name hellosvc2 \
--network msdemonet \
duncanandrew/hellosvc

docker run -d \
--name gateway \
--network msdemonet \
-p 443:443 \
-v $(pwd)/gateway/tls:/tls:ro \
-e CERTPATH=/tls/fullchain.pem \
-e KEYPATH=/tls/privkey.pem \
-e HELLOSVCADDR=hellosvc1,hellosvc2 \
duncanandrew/gateway
