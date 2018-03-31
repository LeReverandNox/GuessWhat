#!/bin/bash
if [ ${ENVIRONMENT} == "prod" ]; then
	exec ./guesswhat
else
	exec fresh
fi
