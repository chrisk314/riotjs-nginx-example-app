#!/bin/bash

set -euo pipefail

while getopts "d:e:s" opt; do
  case ${opt} in
    d) domains+=("${OPTARG}");;
    e) email=${OPTARG};;
    s) staging=1
  esac
done
shift $(( OPTIND - 1 ))

staging=${staging:-0}
rsa_key_size=4096

if [ -z "${email+set}" ] || [ -z "${email+unset}" ]; then
  echo "Missing required argument: -e <email>"; exit 1;
fi
if [ -z "${domains+set}" ] || [ -z "${domains+unset}" ]; then
  echo "Missing required argument: -d <domain_1> [..., -d <domain_N>]"; exit 1;
fi

echo "### Creating certbot docker volumes and helper container ..."
docker volume create certbot-etc
docker volume create certbot-www
docker run -d --rm --name helper \
  -v certbot-etc:/etc/letsencrypt -v certbot-www:/var/www/certbot \
  --entrypoint "/bin/sh" certbot/certbot -c "tail -f /dev/null"
echo

echo "### Downloading recommended TLS parameters ..."
docker exec helper sh -c "\
  wget https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf \
    -O /etc/letsencrypt/options-ssl-nginx.conf"
docker exec helper sh -c "\
  wget https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem \
    -O /etc/letsencrypt/ssl-dhparams.pem"
echo

echo "### Creating dummy certificate for ${domains[0]} ..."
path="/etc/letsencrypt/live/${domains[0]}"
docker exec helper sh -c "\
  mkdir -p ${path} && \
  openssl req -x509 -nodes -newkey rsa:1024 -days 1\
    -keyout '${path}/privkey.pem' \
    -out '${path}/fullchain.pem' \
    -subj '/CN=localhost'"
echo

echo "### Starting nginx ..."
docker-compose -f docker-compose.yml up --force-recreate -d nginx
echo

echo "### Deleting dummy certificate for ${domains[0]} ..."
docker exec helper sh -c "\
  rm -Rf /etc/letsencrypt/live/${domains[0]} && \
  rm -Rf /etc/letsencrypt/archive/${domains[0]} && \
  rm -Rf /etc/letsencrypt/renewal/${domains[0]}.conf"
echo

echo "### Requesting Let's Encrypt certificate for ${domains[@]} ..."
#Join $domains to -d args
domain_args=""
for domain in "${domains[@]}"; do
  domain_args="${domain_args} -d ${domain}"
done

# Select appropriate email arg
case "${email}" in
  "") email_arg="--register-unsafely-without-email" ;;
  *) email_arg="--email ${email}" ;;
esac

# Enable staging mode if needed
staging_arg=$((( ${staging} )) && echo "--staging" || echo)

docker exec helper sh -c "\
  certbot certonly --webroot -w /var/www/certbot \
    ${staging_arg} ${email_arg} ${domain_args} \
    --rsa-key-size ${rsa_key_size} \
    --agree-tos \
    --force-renewal"
echo

echo "### Tearing down helper and nginx ..."
docker stop helper &
docker-compose -f docker-compose.yml down &

wait

exit 0
