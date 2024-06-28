run:
  docker run --rm \
    -v $(pwd)/conf/Server.xml:/opt/ovenmediaengine/bin/origin_conf/Server.xml \
    -p 1935:1935 -p 9999:9999/udp -p 9000:9000 -p 3333:3333 \
    -p 3478:3478 -p 10000-10009:10000-10009/udp \
    airensoft/ovenmediaengine:0.16.4

ingest:
    ffmpeg -re -f lavfi -i testsrc=size=1280x720:rate=30 \
        -vf "drawtext=fontsize=30:fontcolor=white:x=7:y=7:text='Time\: %{localtime\:%X}',format=yuv420p" \
        -pix_fmt yuv420p -c:v libx264 -preset ultrafast \
        -b:v 600k -max_muxing_queue_size 1024 -tune zerolatency \
        -g 30 -f flv rtmp://localhost:1935/app/lowlatency
