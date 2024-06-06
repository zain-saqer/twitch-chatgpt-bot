#!make
include .env
export $(shell sed 's/=.*//' .env)

test:
	env
build:
	@docker build -t ${APP_IMAGE} -f ./docker/app/Dockerfile .

push-image:
	@docker image push ${APP_IMAGE}

pull-image:
	@docker image pull ${APP_IMAGE}

up:
	@docker stack deploy --compose-file=docker-stack.yml twitch-chatgpt-bot

down:
	@docker stack down twitch-chatgpt-bot

app-service-logs:
	@docker service logs twitch-chatgpt-bot_app