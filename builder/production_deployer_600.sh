#!/bin/bash

  chmod 774 deployer.sh
	
  if [ -z "$1" ]
  then 	
    echo "Not specified binary folder"
    exit 2
  fi

  if [ -z "$2" ]
  then 	
    echo "Not specified target release folder"
    exit 2
  fi

  OUT_PATH=$1
  RELEASE_FOLDER=$2
  REMOTE_HOST_SV3="sv3"  


  ./deployer.sh ${OUT_PATH} device-queclink_gv600 ${REMOTE_HOST_SV3} ${RELEASE_FOLDER}

  if [ $? != 0 ]
  then
    echo "Failed deploy device-queclink gv600 to ${REMOTE_HOST_SV3}"	
    exit 1 
  fi
