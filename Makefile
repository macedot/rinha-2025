# https://github.com/andersongomes001/rinha-2025

DOCKER_USER := "macedot"
APP_NAME := "rinha-2025"
IMAGE_NAME := "$(DOCKER_USER)/$(APP_NAME)"
MY_VAR := $(git rev-parse --short HEAD)

dev:
	@echo docker compose -f ./docker-compose-dev.yml down && docker compose -f ./docker-compose-dev.yml  build && docker compose -f docker-compose-dev.yml up

build:
	@echo "🐳 Build da imagem Docker..."
	@echo docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .

push: build
	@echo "🔐 Enviando imagens..."
	@echo docker push $(IMAGE_NAME):$(VERSION)
	@echo docker push $(IMAGE_NAME):latest

prod: push
	@echo docker compose down && docker compose up

.PHONY: dev build push prod
