#!/bin/bash

set -e
# Поднимается в корень backend-проекта из паки script

# Backend and Frontend размещаются в одной директории

#backend
backend_image_name="gostudy-go-backend"
backend_container_name="go-backend"
docker_file_path="."

# Stop and remove backend container
if docker ps -q -f "name=${backend_container_name}" | grep -q .; then
  echo "Stopping and removing container: $backend_container_name"
  docker stop "$backend_container_name"
  docker rm "$backend_container_name"
else
  echo "No container with name '$backend_container_name' found."
fi

# Remove the old backend image
echo "Removing image: $backend_image_name"
docker rmi "$backend_image_name" || true


docker compose -f docker-compose.yaml -p gostudy up -d