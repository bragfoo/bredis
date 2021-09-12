#! /bin/bash
ab -c 4 -n 1000 127.0.0.1:8080/set/b/redis
ab -c 4 -n 1000 127.0.0.1:8080/get/b

# set diffent val
siege -f ./requests
