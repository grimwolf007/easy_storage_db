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
    file="./test/test_healthcheck_files/healthcheck_testing"
    test_file="${file}.test"
    actual_file="${file}.act"
    count=1
    until wget 127.0.0.1:8080/health-check -O $actual_file  || [ $count -eq 10 ]
    do
      let count+=1
      sleep 1
    done
    diff $test_file $actual_file && tmp=0 || tmp=1
    SucOrFail $tmp "Health_Check_${count} Test"
  echo

  info "Testing /upload/server"
    file="./test/test_upload_files/upload_testing"
    test_file="${file}.test"
    upload_file="./test/upload_testing.test"
    actual_file="${file}.act"
    curl -f -X POST http://localhost:8080/upload/server -F "upload[]=@${test_file}"
    mv -f ${upload_file} ${actual_file}
    diff $test_file $actual_file && tmp=0 || tmp=1
    SucOrFail $tmp "Upload Test"
  echo

info Testing Minio
  test="List_Minio_Buckets"
  info "Testing ${test}"
    file="./test/test_list-buckets_files/list-buckets_testing"
    test_file="${file}.test"
    actual_file="${file}.act"
    curl -o $actual_file -f -X GET http://localhost:8080/list_buckets
    exists=$(cat $actual_file | jq '.[]')
    info "Buckets: $exists"
    [ ! -z "$exists" ] && tmp=0 || tmp=1
    SucOrFail $tmp "$test Test"
  echo

  test="Create_Minio_Bucket"
  info "Testing ${test}"
    file="./test/test_make-bucket_files/make-bucket_testing"
    test_file="${file}.test"
    actual_file="${file}.act"
    curl -s -f -X POST http://localhost:8080/create_bucket/test-bucket2
    curl -o $actual_file -f -X GET http://localhost:8080/list_buckets
    exists=$(cat $actual_file | jq '.[] | select(.name=="test-bucket2")')
    info "New Bucket: $exists"
    [ ! -z "$exists" ] && tmp=0 || tmp=1
    SucOrFail $tmp "$test Test"
  echo

  test="Remove_Minio_Bucket"
  info "Testing ${test}"
    file="./test/test_remove-bucket_files/remove-bucket_testing"
    test_file="${file}.test"
    actual_file="${file}.act"
    curl -s -f -X POST http://localhost:8080/remove_bucket/test-bucket2
    curl -o $actual_file -f -X GET http://localhost:8080/list_buckets
    exists=$(cat $actual_file | jq '.[] | select(.name=="test-bucket2")')
    
    [ -z "$exists" ] && tmp=0 || tmp=1
    [ -z "$exists" ] && info "test-bucket2 nolonger exists."
    [ ! -z "$exists" ] && fatal "$exists"
    SucOrFail $tmp "$test Test"
  echo

  test="Get_Objects"
  info "Testing ${test}"
    file="./test/test_${test}_files/${test}_testing"
    test_file="${file}.test"
    actual_file="${file}.act"
    curl -o $actual_file -f -X GET http://localhost:8080/bucket/test-bucket
    exists=$(cat $actual_file | jq '.[0]')
    info "Objects: $exists"
    [ ! -z "$exists" ] && tmp=0 || tmp=1
    SucOrFail $tmp "$test Test"
  echo

  test="Put_Object"
  info "Testing ${test}"
    name="upload_testing2"
    file="./test/test_upload_files/${name}"
    test_file="${file}.test"
    upload_file="./test/${name}.test"
    actual_file="${file}.act"
    curl -X POST http://localhost:8080/upload/object/test-bucket \
       -F "upload[]=@${test_file}"
    curl -o test.json -f -X GET http://localhost:8080/bucket/test-bucket
    exists=$(cat test.json | jq '.[1]')
    info "Objects: $exists"
    [ ! -z "$exists" ] && tmp=0 || tmp=1
    SucOrFail $tmp "${test} Test"
    rm test.json
  echo 



  test="Remove_Object"
  info "Testing ${test}"
    curl -f -X POST http://localhost:8080/remove_object/test-bucket/upload-testing2.test
    curl -o test.json -f -X GET http://localhost:8080/bucket/test-bucket
    exists=$(cat test.json | jq '.[1]')
    info "Objects: $exists"
    [ -z "$exists" ] && tmp=0 || tmp=1
    SucOrFail $tmp "${test} Test"
    rm test.json
  echo

#echo Testing Postgres

#echo Testing Adminer






info "Pausing"
read
echo Stopping server
docker compose down
