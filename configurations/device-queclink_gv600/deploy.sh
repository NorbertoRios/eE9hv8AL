#!/bin/sh
CURRENT=current
LOG_DIR=logs
SERVICE_NAME=queclink600
SERVICE_HOME=~/device-queclink_gv600

echo "----------------------------Starting deploy------------------------------"


echo "---------------------Runing docker service-------------------------------"
docker ps -a | head -1
docker ps -a | grep ${SERVICE_NAME}

echo "-------------------------------------------------------------------------"

cd ${SERVICE_HOME}

if [ -e ${CURRENT} ]
then
        rm -r ${CURRENT}
        echo "Removed \"${CURRENT}\" symlink";
fi
sleep 1

ln -fs $1 ${CURRENT}
echo "Created new symlink \"${CURRENT}\" on \"$1\""
sleep 1

if [ ! -d config ];
then
  mkdir config
  echo "Created directory config"
fi


#tar existing configurations
if [ ! -d config/backup ];
then
  mkdir config/backup
  echo "Created directory config/backup"
fi

tar -zcvf config/backup/config.$(date +%Y-%m-%d_%H-%M-%S).tar.gz  config/*.xml config/*.json
cp -fr ${CURRENT}/configurations/*.* config/

if [ ! -d ${LOG_DIR} ]; then
  mkdir ${LOG_DIR}
  echo "Created directory \"/${LOG_DIR}\""
fi

echo "-----------------------------------Building docker image---------------------------------------"
docker-compose build

echo "-----------------------------------Starting docker image---------------------------------------"
docker-compose up -d

echo "-----------------------------------Runing docker container-------------------------------------"
docker ps -a | head -1
docker ps -a | grep ${SERVICE_NAME}
echo "-----------------------------------------------------------------------------------------------"
exit 0

