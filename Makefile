DIR=deployments/docker
RECIPE=${DIR}/docker-compose.yaml
NAMESPACE=builder${COMPONENT}
MIGRATION_PATH=`pwd`/migrations/test

include .env
export $(shell sed 's/=.*//' .env)


DIND_PREFIX ?= $(HOME)
ifneq ($(HOST_PATH),)
DIND_PREFIX := $(HOST_PATH)
endif
ifeq ($(CACHE_PREFIX),)
	CACHE_PREFIX=/tmp
endif

PREFIX=$(shell echo $(PWD) | sed -e s:$(HOME):$(DIND_PREFIX):)
UID=$(shell whoami)

IMAGE_TAG ?= master
export $IMAGE_TAG

.PHONY : test

test: infratest
	docker run \
		--network ${NAMESPACE}_default \
		--env-file .env \
		-v $(CACHE_PREFIX)/cache/go:/go/pkg/mod \
		-v $(CACHE_PREFIX)/cache/apk:/etc/apk/cache \
		-v $(PREFIX)/deployments/docker/build-results:/build \
		-v $(PREFIX)/:/src \
		-v $(PREFIX)/migrations:/migrations \
		-v $(PREFIX)/scripts/test.sh:/test.sh \
		-e UID=$(UID) \
		golang:1.12-alpine /test.sh 

deps:
	@./scripts/deps.sh
	
cleantest:
	docker-compose -f ${RECIPE} -p ${NAMESPACE} stop 
	docker-compose -f ${RECIPE} -p ${NAMESPACE} rm -f testnats
	docker network remove ${NAMESPACE}_default ; /bin/true 

infratest:
	docker network create -d bridge ${NAMESPACE}_default ; /bin/true 
	docker-compose -f ${RECIPE} -p ${NAMESPACE} up -d --force-recreate testnats

build:
	docker run -v $(CACHE_PREFIX)/cache/go:/go/pkg/mod \
		-v $(CACHE_PREFIX)/cache/apk:/etc/apk/cache \
		-v $(PREFIX)/deployments/docker/build-results:/build \
		-v $(PREFIX)/scripts/build.sh:/build.sh \
		-v $(PREFIX)/:/src \
		-v $(PREFIX)/cmd:/src/cmd \
		-e RUNNER_USER \
		golang:1.12-alpine /build.sh ${COMPONENT}

clean:
	docker-compose -f ${RECIPE} -p ${NAMESPACE} stop 
	docker-compose -f ${RECIPE} -p ${NAMESPACE} rm -f exporter

exporter: clean build 
	docker-compose -f ${RECIPE} -p ${NAMESPACE} build --no-cache exporter

