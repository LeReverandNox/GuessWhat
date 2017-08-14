#!/bin/bash
if [ ${ENVIRONMENT} == "prod" ]; then
	app
else
	fresh
fi