BUILD_VERSION   := $(shell cat version)
BUILD_DATE      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)

all:
	./cross_compile.sh

release: clean all
	ghr -u mritd -t ${GITHUB_TOKEN} -replace -recreate -name "Bump ${BUILD_VERSION}" --debug ${BUILD_VERSION} dist

docker:
	docker build -t ytpay/telepush:${BUILD_VERSION} .

clean:
	rm -rf dist

install:
	go install

.PHONY : all release clean install
