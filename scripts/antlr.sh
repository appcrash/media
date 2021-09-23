#!/bin/bash

# get anltr by:
# curl -O https://www.antlr.org/download/antlr-4.9.2-complete.jar


ANTLR_JAR=/usr/local/lib/antlr-4.9.2-complete.jar
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

export CLASSPATH=".:$ANTLR_JAR:$CLASSPATH"
export antlr4="java -jar $ANTLR_JAR "

cd $DIR/../server/comp/nmd && $antlr4 -Dlanguage=Go -package nmd nmd.g4 
