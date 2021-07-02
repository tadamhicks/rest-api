#!/bin/bash
docker run -d -p 8000:8000 -e MONGO_DATABASE='person_db' -e MONGO_USERNAME='' -e MONGO_PASSWORD='' -e MONGO_SERVER='192.168.1.108' tadamhicks/resting-test
