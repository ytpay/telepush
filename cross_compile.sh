#!/usr/bin/env bash

set -e

GO111MODULE="on"
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64"

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/telepush_${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/tgsend_${GOOS}_${GOARCH}.exe
    fi

    echo "build => ${TARGET}"
    go build -trimpath -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"
done