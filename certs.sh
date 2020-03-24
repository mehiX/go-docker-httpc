#!/bin/bash

CERTS_PATH=certs

if [ -d "${CERTS_PATH}" ]; then
  echo Directory ${CERTS_PATH} already exists
  exit 1
else 
  mkdir -p ${CERTS_PATH}
fi

openssl genrsa -out ${CERTS_PATH}/key.pem

openssl req -new -key ${CERTS_PATH}/key.pem -out ${CERTS_PATH}/cert.pem -nodes -subj "/CN=localhost"

openssl req -x509 -days 20 -key ${CERTS_PATH}/key.pem -in ${CERTS_PATH}/cert.pem -out ${CERTS_PATH}/certificate.pem

echo Certificate created at path: ${CERTS_PATH}