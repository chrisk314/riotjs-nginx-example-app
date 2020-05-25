#!/usr/bin/env bash
SSL_DIR=${HOME}/.ssl
mkdir -p ${SSL_DIR}/{certs,private}
openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 \
  -keyout ${SSL_DIR}/private/RootCA.key -out ${SSL_DIR}/certs/RootCA.pem -subj "/C=US/CN=Localhost-Root-CA"
openssl x509 -outform pem -in ${SSL_DIR}/certs/RootCA.pem -out ${SSL_DIR}/certs/RootCA.crt
openssl req -new -nodes -newkey rsa:2048 -keyout ${SSL_DIR}/private/localhost.key \
  -out ${SSL_DIR}/certs/localhost.csr -subj "/C=UK/ST=London/L=London/O=Localhost-Certificates/CN=localhost.local"
cat <<EOF > ${SSL_DIR}/domains.ext
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = localhost.lan
EOF
openssl x509 -req -sha256 -days 1024 -in ${SSL_DIR}/certs/localhost.csr -CA ${SSL_DIR}/certs/RootCA.pem \
  -CAkey ${SSL_DIR}/private/RootCA.key -CAcreateserial -extfile ${SSL_DIR}/domains.ext -out ${SSL_DIR}/certs/localhost.crt
openssl x509 -inform PEM -in ${SSL_DIR}/certs/localhost.crt > ${SSL_DIR}/certs/localhost.pem
openssl rsa -in ${SSL_DIR}/private/localhost.key -text > ${SSL_DIR}/private/localhost.pem
cat ${SSL_DIR}/certs/localhost.crt ${SSL_DIR}/private/localhost.key > ${SSL_DIR}/certs/localhost.includesprivatekey.pem
chmod 0644 ${SSL_DIR}/certs/*
chmod 0640 ${SSL_DIR}/private/*

