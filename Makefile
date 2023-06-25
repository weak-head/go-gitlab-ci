# Docker Image
IMAGE_ID ?= gogin:local

.PHONY: docker
docker:
	@docker build -t ${IMAGE_ID} .

.PHONY: run
run:
	@docker run --rm \
		-p "8080:8080" \
		${IMAGE_ID}

.PHONY: install
install:
	helm install \
		--namespace=services \
		--create-namespace \
		gogin deploy/gogin

.PHONY: uninstall
uninstall:
	helm uninstall \
		--namespace=services \
		gogin