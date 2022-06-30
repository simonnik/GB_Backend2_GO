#write here your username
USERNAME := sim4nix
APP_NAME := k8s-go-app
VERSION := latest

#write here path for your project
PROJECT := k8s-go-app
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := latest
APP_NAME := k8s-go-app

all: run

run:
	go install -ldflags="-X '$(PROJECT)/version.Version=$(VERSION)' \
	-X '$(PROJECT)/version.Commit=$(GIT_COMMIT)'" && $(APP_NAME)

build_container:
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) \
	--build-arg=VERSION=$(VERSION) --build-arg=PROJECT=$(PROJECT) \
	-t docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) .

run_container:
	docker run -p 8080:8080 docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)

push_container:
	docker push docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)

apply_deploy:
	kubectl apply -f .kube/deployment.yaml

apply_service:
	kubectl apply -f .kube/service.yaml

apply_ingress:
	kubectl apply -f .kube/ingress.yaml

