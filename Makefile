
.PHONY: install uninstall

install:
	helm install \
		--namespace=services \
		--create-namespace \
		gogin helm

uninstall:
	helm uninstall \
		--namespace=services \
		gogin
