#!/bin/bash
set target=""
if [ -f .env ]
then
	export $(cat .env | xargs)
	target=$STTYLUS_PRODUCTION_TARGET
fi

printf "STTylus - Deployment tool\n"
printf "Working directory $(pwd) \n"
printf "target: $target\n"
printf "==============================\n\n\n"
printf "Retreiving backup of sttylus.db -> production.db\n"
scp $user@$target/sttylus.db production.db
printf "\n"
printf "Pushing master to production\n"
git push production master
