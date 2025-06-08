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

set -e
set -o pipefail

token=$(curl  -fsSL "${baseurl}/auth/firsttime" -H "Content-Type: application/json" --data-raw "${firsttime_payload}" | jq -r '.token')
org_id=$(curl  -fsSL "${baseurl}/user"           -H "Authorization: Bearer ${token}" | jq -r '.organizations[0].id')
apikey=$(curl -X POST -fsSL "${baseurl}/keys"           -H "Content-Type: application/json" -H "X-Organization: ${org_id}" -H "Authorization: Bearer ${token}"  --data-raw "${createkey_payload}" | jq -r '.key.key')

echo ${apikey}
