#!/bin/bash
old=`cut -d '=' -f2 .env.build`
incr=`expr $old + 1`
sed -i "s/STTYLUS_BUILD=$old\$/STTYLUS_BUILD=$incr/" .env.build
