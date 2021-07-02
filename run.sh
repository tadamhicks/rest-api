#!/bin/bash
docker run -d -p 8000:8000 \
-e MONGO_DATABASE='person_db' \
-e MONGO_USERNAME='' \
-e MONGO_PASSWORD='' \
-e MONGO_SERVER='192.168.1.108' \
-e HONEYCOMB_APIKEY='c4b05d6b2259d9d6fca768d4ba9c811a' \
-e HONEYCOMB_DATASET='restful-sleep' \
-e HONEYCOMB_SERVICENAME='restful-sleep-svc' \
tadamhicks/resting-test
