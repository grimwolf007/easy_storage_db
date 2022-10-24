echo Adding to Build Archive
  go build -o builds/webserver_`date +%F.%T`
echo Replacing Latest
  go build -o latest_build/webserver
echo Building Docker compose
  docker compose build
echo Starting server
  docker compose up -d
  sleep 10
echo Testing webapp
  echo Testing /ping 
    wget 127.0.0.1:8080/ping -O ./test/tmp
    if diff ./test/tmp ./test/ping.test;
      then
        echo "Ping Test Successful"
      else
        echo "Ping Test Failed"
    fi
    
echo Testing Minio

echo Testing Postgres

echo Testing Adminer








echo Stopping server
  docker compose down
