#!/bin/bash

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "Missing required positional argument(s): <domain_1> [..., <domain_N>]"; exit 1;
else
  domains=("$@")
fi

nginx_conf_path="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"/../nginx/prod.nginx.conf

sed -i "s!example.com/fullchain.pem!${domains[0]}/fullchain.pem!g" ${nginx_conf_path}
sed -i "s!example.com/privkey.pem!${domains[0]}/privkey.pem!g" ${nginx_conf_path}
sed -i "s!server_name example.com!server_name $(echo ${domains[@]})!g" ${nginx_conf_path}

exit 0
