# Building Docker images for ARM on macOS

I'm using the experimental Buildx support in Docker Desktop for Mac to make it easier to build the base images for my containers that will run on the Raspberry Pi.

### Example

On the Mac:

```sh
cd images/ubuntu-bluetooth
docker buildx build --platform linux/arm64 -t mmoriarity/ubuntu-bluetooth:latest --push .
```

On the Raspberry Pi:

```sh
sudo docker run --rm -it --net=host mmoriarity/ubuntu-bluetooth /bin/bash
```

Tada!
