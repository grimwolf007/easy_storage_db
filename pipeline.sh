#!/bin/bash

#Colors
RESET='\033[0m'
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
LRED='\033[1;31m'

info(){
       	echo -e "${CYAN}${@}${RESET}"
}

pass(){
       	echo -e "${GREEN}${@}${RESET}"
}

warn(){
       	echo -e "${YELLOW}${@}${RESET}"
}

fail(){
       	echo -e "${LRED}${@}${RESET}"
}

fatal(){
       	echo -e "${RED}${@}${RESET}"
}

SucOrFail(){
	b=$1
	shift
	name="$@"

	if [ $b -eq 0 ]
       	then	
          pass "${name} Successful"
	else
	  fail "${name} Failed"
	fi

}

info Adding to Build Archive
  go build -o builds/webserver_`date +%F.%T`
info Replacing Latest
  go build -o latest_build/webserver
info Building Docker compose
  docker compose build
info Starting servers
  docker compose up &

  
if [ "$1" == "prune" ]
then
  info Removing unused images/containers
    docker system prune -f
fi


info "Testing webapp"

  info "Testing /healthcheck"
    count=1
    until wget 127.0.0.1:8080/health-check -O ./test/tmp  || [ $count -eq 10 ]
    do
      let count+=1
      sleep 1
    done
    diff ./test/tmp ./test/healthcheck.test && tmp=0 || tmp=1
    SucOrFail $tmp "Health_Check_${count} Test"
    rm -f ./test/tmp
  echo

  info "Testing /upload"
    test_file="/home/john/git/easy_storage_db/test/test_upload_files/testing.txt"
    check_file="/home/john/git/easy_storage_db/test/testing.txt"
    curl -s -f -X POST http://localhost:8080/upload \
      -F "upload[]=@${test_file}"
    diff $test_file $check_file && tmp=0 || tmp=1
    SucOrFail $tmp "Upload Test"
    rm -f $check_file
  echo

info Testing Minio

  info "Testing get minio bucket"
  curl -s -f -X GET http://localhost:8080/get_bucket/test-bucket && tmp=0 || tmp=1
  SucOrFail $tmp "Get_minio_test_bucket Test"
  echo

  info "Testing make minio bucket"
  curl -s -f -X POST http://localhost:8080/create_bucket/test-bucket2 && tmp=0 || tmp=1
  SucOrFail $tmp "Create_minio_test_bucket Test"
  echo

  info "Testing list minio buckets"
  fatal "No Test Made"
  echo

  info "Testing remove minio bucket"
  fatal "No Test Made"
  echo

  info "Testing add new object"
  fatal "No Test Made"
  echo 

  info "Testing get object"
  fatal "No Test Made"
  echo

  info "Testing remove object"
  fatal "No Test Made"
  echo

#echo Testing Postgres

#echo Testing Adminer






info "Pausing"
read
echo Stopping server
docker compose down
