TIMESTAMP=$(shell date +"%Y%m%d%H%M")

default:
	@echo "docker buildx create"
	DOCKER_BUILDKIT=1 docker buildx create --use

	DOCKER_BUILDKIT=1 docker buildx build --platform linux/386,linux/amd64,linux/arm64/v8 -t registry.sean.farm/mattermost-notify:latest . --push
	