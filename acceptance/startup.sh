#!/bin/bash

baseurl="http://localhost:3100"

read -r -d '' firsttime_payload <<EOS
{
  "companyname": "testaccount",
  "name": "testaccount",
  "email": "testaccount@nomail.com",
  "password": "testaccount"
}
EOS

read -r -d '' createkey_payload <<EOS
{
  "description":"testkey",
  "type":"user"
}
EOS

retry_curl() {
  local out

  while true; do
    if out=$(curl -fsSL "$@" 2>/dev/null); then
      if [ -n "$out" ]; then
        printf '%s\n' "$out"
        return 0
      fi
    fi

    echo "curl failed or empty response, retrying in 1s..." >&2
    sleep 1
  done
}

set -e
set -o pipefail

token=$(
  retry_curl "${baseurl}/auth/firsttime" \
    -H "Content-Type: application/json" \
    --data-raw "${firsttime_payload}" \
  | jq -r '.token'
)

org_id=$(
  retry_curl "${baseurl}/user" \
    -H "Authorization: Bearer ${token}" \
  | jq -r '.organizations[0].id'
)

apikey=$(
  retry_curl -X POST "${baseurl}/keys" \
    -H "Content-Type: application/json" \
    -H "X-Organization: ${org_id}" \
    -H "Authorization: Bearer ${token}" \
    --data-raw "${createkey_payload}" \
  | jq -r '.key.key'
)

echo "${apikey}"
