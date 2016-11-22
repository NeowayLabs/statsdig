#!/bin/bash

set -o nounset
set -o errexit

TAGS="role:webserver,location:europe"

docker run -d --name sysdig-agent --privileged --net host --pid host -e ACCESS_KEY=$SYSDIG_ACCESS_KEY -e TAGS=$TAGS -v /var/run/docker.sock:/host/var/run/docker.sock -v /dev:/host/dev -v /proc:/host/proc:ro -v /boot:/host/boot:ro -v /lib/modules:/host/lib/modules:ro -v /usr:/host/usr:ro sysdig/agent
