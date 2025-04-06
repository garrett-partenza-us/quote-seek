docker buildx build --platform linux/arm64 -t ghcr.io/garrett-partenza-us/quote-seek/backend:latest -f ./Dockerfile.backend --push .
docker buildx build --platform linux/arm64 -t ghcr.io/garrett-partenza-us/quote-seek/frontend:latest -f ./Dockerfile.frontend --push .
