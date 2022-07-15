
.PHONY: build
build:
	docker build -t registry.lothric.net/gogin:0.0.1 .

.PHONY: push
push: build
	docker push registry.lothric.net/gogin:0.0.1

.PHONY: install
install:
	helm install --namespace=default gogin deploy/gogin

.PHONY: uninstall
uninstall:
	helm uninstall --namespace=default gogin