for c in $(docker ps -aq); do docker rm -f "$c"; done
docker system prune -a --volumes -f