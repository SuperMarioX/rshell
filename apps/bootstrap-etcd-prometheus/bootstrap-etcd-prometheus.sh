#!/bin/env bash
scriptFilePath="$(dirname ${BASH_SOURCE[0]})"
sourceFilePath="$(readlink -f ${scriptFilePath})"
sourceFileName="$(basename ${BASH_SOURCE[0]})"

LOGFILE=${sourceFilePath}/${sourceFileName}.log
PIDFILE=${sourceFilePath}/${sourceFileName}.pid

echo "Begin bootstrap ectd prometheus. $(date)" | tee -a ${LOGFILE}

[ $# -ne 2 ] && { echo 'arg number not == 2' | tee -a ${LOGFILE}; exit 1; }

DATA_DIR=$1
mkdir -p ${DATA_DIR}

CONFIG_FILE=$2
[ ! -e "${CONFIG_FILE}" ] && { echo 'config file not found' | tee -a ${LOGFILE}; exit 1; }

OPTIONS="--config.file ${CONFIG_FILE} \
  --web.listen-address=0.0.0.0:9090 \
  --storage.tsdb.path ${DATA_DIR} \
  --storage.tsdb.retention=15d"

echo "options: ${OPTIONS}" >> ${LOGFILE}
	
netstat -anop | grep -E ":9090 " >> ${LOGFILE}
[ $? -eq 0 ] && { echo "9090 in use now" | tee -a ${LOGFILE}; exit 1; }

[ ! -e "${sourceFilePath}/prometheus" ] && { echo 'prometheus not found' | tee -a ${LOGFILE}; exit 1; }
chmod +x ${sourceFilePath}/prometheus
${sourceFilePath}/prometheus ${OPTIONS} >> ${LOGFILE} 2>&1 &

ps -ef | grep "${sourceFilePath}/prometheus" | grep -v grep >> ${LOGFILE}

echo "Bootstrap etcd prometheus finish" | tee -a ${LOGFILE}

