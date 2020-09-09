#!/bin/bash


	  if [ -z "$1" ]
	  then 	
	    echo "Not specified target deploy configuration"
	    exit 2
	  fi

	  if [ -z "$2" ]
	  then 	
	    echo "Not specified build files output"
	    exit 2
	  fi

	  if [ -z "$3" ]
	  then 	
	    echo "Not specified target server"
	    exit 2
	  fi
	  if [ -z "$4" ]
	  then 	
	    echo "Not specified release folder"
	    exit 2
	  fi


OUT_PATH=$1
DEPLOY_DESTINATION=$2
REMOTE=$3	
RELEASE_FOLDER=$4

	  echo "Starting upload to '"${REMOTE}"'"

	  ssh ${REMOTE} mkdir ${DEPLOY_DESTINATION}
	  rsync -avz ../configurations/${DEPLOY_DESTINATION}/docker/docker-compose.yml ${REMOTE}:${DEPLOY_DESTINATION}/
	  rsync -avz ../configurations/${DEPLOY_DESTINATION}/deploy.sh ${REMOTE}:${DEPLOY_DESTINATION}/
	  rsync -avz ../configurations/${DEPLOY_DESTINATION}/docker/Dockerfile ${REMOTE}:${DEPLOY_DESTINATION}/${RELEASE_FOLDER}/
	  rsync -avz ../configurations/${DEPLOY_DESTINATION}/config/credentials.json ${REMOTE}:${DEPLOY_DESTINATION}/${RELEASE_FOLDER}/configurations/
	  rsync -avz ../configurations/${DEPLOY_DESTINATION}/config/ReportConfiguration.xml ${REMOTE}:${DEPLOY_DESTINATION}/${RELEASE_FOLDER}/configurations/

	  rsync -avz ${OUT_PATH}/src  ${REMOTE}:${DEPLOY_DESTINATION}/${RELEASE_FOLDER}/
	  rsync -avz ${OUT_PATH}/vendor ${REMOTE}:${DEPLOY_DESTINATION}/${RELEASE_FOLDER}/
	  ssh ${REMOTE} chmod 774 ${DEPLOY_DESTINATION}/deploy.sh
	  ssh ${REMOTE} . ${DEPLOY_DESTINATION}/deploy.sh ${RELEASE_FOLDER} ${DEPLOY_DESTINATION}
	  echo "Deploy ${DEPLOY_DESTINATION} to ${REMOTE} complete"
    exit $?