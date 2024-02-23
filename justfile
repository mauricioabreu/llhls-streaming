run:
  docker run --name ome \
    -v $(pwd)/conf/Server.xml:/opt/ovenmediaengine/bin/origin_conf/Server.xml \
    -p 1935:1935 -p 9999:9999/udp -p 9000:9000 -p 3333:3333 \
    -p 3478:3478 -p 10000-10009:10000-10009/udp \
    airensoft/ovenmediaengine:0.16.4
