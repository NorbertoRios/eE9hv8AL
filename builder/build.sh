#!/bin/bash
GIT_REVISION=$(git describe --long --tags)
CURR_FOLDER=$(pwd)
RELEASE_FOLDER='queclink_'${GIT_REVISION}
OUT_PATH=${CURR_FOLDER}/../
	case "$1" in
	 deploy)
	  chmod 774 production_deployer.sh 
          ./production_deployer.sh ${OUT_PATH} ${RELEASE_FOLDER}
	 ;;
	deploy600)
	  chmod 774 production_deployer_600.sh 
          ./production_deployer_600.sh ${OUT_PATH} ${RELEASE_FOLDER}
	 ;;
	 usage)
	  echo "Usage of build script: build.sh [action]"
	  echo "	actions: build, deploy, usage"
	 ;;
	esac
exit $?
