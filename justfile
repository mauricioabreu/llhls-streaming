run:
  docker compose up

stop:
  docker compose down

ingest-srt:
  ffmpeg -re -f lavfi -i testsrc=size=1280x720:rate=30 \
    -vf "drawtext=fontsize=30:fontcolor=white:x=7:y=7:text='Time\: %{localtime\:%X}',format=yuv420p" \
    -pix_fmt yuv420p -c:v libx264 -preset ultrafast \
    -b:v 600k -max_muxing_queue_size 1024 -tune zerolatency \
    -g 30 -f mpegts "srt://127.0.0.1:9999?streamid=srt://127.0.0.1:9999/ll/llhls"
