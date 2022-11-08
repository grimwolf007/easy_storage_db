#!/bin/bash
log_rotate(){
  l="../latest_pipeline.log"
  swp="./swp.log"
  number=$1
  let "number-=1"
  arr=$(seq 1 $number)
  cd ./pipeline_logs
  touch $l
  touch $swp
  for i in $arr
    do  
      file=pipeline_$i.log
      touch $file
      cp $file $swp
      cp $l $file
      cp $swp $l
  done
  cd ..
}

log_rotate 10
date > latest_pipeline.log
./pipeline.sh $1  2>&1 | tee latest_pipeline.log
less -R latest_pipeline.log
