#!/bin/bash
json='{"private_key":"'"$2"'","vote_id":"'"$3"'","codes":['
for ((i=1; i<=$4; i++))
do
    if [ $i != 1 ]
    then
        json=$json','
    fi
    json=$json'"'$i'"'
done
json=$json']}'
curl -X POST -H "Content-Type: application/json" -d $json "$1"/ballot/add