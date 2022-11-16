#!/bin/bash
set target=""
set user=""
if [ -f .env ]
then
	export $(cat .env | xargs)
	target=$STTYLUS_PRODUCTION_TARGET
  user=$STTYLUS_PRODUCTION_USER
fi

printf "STTylus - Deployment tool\n"
printf "Working directory $(pwd) \n"
printf "target: $target\n"
printf "==============================\n\n\n"
printf "Retreiving backup of sttylus.db -> production.db\n"
rsync -auv $user@$target/sttylus.db production.db
printf "\n"

rsync -aunv --exclude-from="exclude.txt" . qwertyist@sttylus.se:beta/
if [ $? -eq 0 ]; then

    echo OK

  read -p "Are you sure? " -n 1 -r
  echo    # (optional) move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    rsync -auv --exclude-from="exclude.txt" . qwertyist@sttylus.se:beta/
  fi
else
    echo FAIL
fi
