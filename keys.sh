#!/bin/bash

KEYS_PATH=keys
PK_PATH=${KEYS_PATH}/app.rsa
PUB_PATH=${PK_PATH}.pub

if [ -d "${KEYS_PATH}" ]; then
    echo "The keys directory (${KEYS_PATH}) already exists. Please use an empty directory"
    exit 1
else
    mkdir -v "${KEYS_PATH}"
fi

echo Generate private key: ${PK_PATH}
openssl genrsa -out ${PK_PATH} 1024

echo Generate public key: ${PUB_PATH}
openssl rsa -in ${PK_PATH} -pubout > ${PUB_PATH}
