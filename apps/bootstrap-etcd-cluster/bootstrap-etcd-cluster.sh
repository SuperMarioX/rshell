#!/bin/env bash
scriptFilePath="$(dirname ${BASH_SOURCE[0]})"
sourceFilePath="$(readlink -f ${scriptFilePath})"
sourceFileName="$(basename ${BASH_SOURCE[0]})"

LOGFILE=${sourceFilePath}/${sourceFileName}.log
PIDFILE=${sourceFilePath}/${sourceFileName}.pid

echo "Begin bootstrap ectd cluster. $(date)" | tee -a ${LOGFILE}

[ $# -ne 4 ] && { echo 'arg number not == 4' | tee -a ${LOGFILE}; exit 1; }

#INITIAL_CLUSTER=infra0=https://10.0.27.239:2380,infra1=https://10.0.20.81:2380,infra2=https://10.0.29.112:2380
INITIAL_CLUSTER=$1
[ -z "${INITIAL_CLUSTER}" ] && { echo 'initial cluster not found' | tee -a ${LOGFILE}; exit 1; }

ROOT_DIR=$2
mkdir -p ${ROOT_DIR}

DATA_DIR=$3
mkdir -p ${DATA_DIR}

INITIAL_CLUSTER_TOKEN=$(echo ${INITIAL_CLUSTER} | md5sum | awk '{print $1}')
INITIAL_CLUSTER_STATE=new

#Generate self-signed certificates by cfssl
#
#echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
#echo '{"signing":{"default":{"expiry":"87600h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
#export ADDRESS=
#export NAME=client
#echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
#export ADDRESS=10.0.27.239,10.0.20.81,10.0.29.112
#export NAME=server
#echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
SSL_DIR=$4
mkdir -p ${SSL_DIR}

[ ! -e ${SSL_DIR}/etcd.ca -o ! -e ${SSL_DIR}/server.cert -o ! -e ${SSL_DIR}/server.key ] && { echo 'server ssl file not found' | tee -a ${LOGFILE}; exit 1; }

NAME=
IP=

oldIFS=$IFS
IFS=,
for host in ${INITIAL_CLUSTER};
do
  n=$(echo ${host} | awk -F= '{print $1}')
  i=$(echo ${host} | awk -F: '{print $2}' | awk -F/ '{print $3}')
  ifconfig -a | grep "inet ${i}" >> ${LOGFILE}
  [ $? -eq 0 ] && { IP=${i}; NAME=${n}; }
done
IFS=$oldIFS

[ -z "${NAME}" -o -z "${IP}" ] && { echo 'ip and name not found' | tee -a ${LOGFILE}; exit 1; }

OPTIONS="--name ${NAME} \
  --data-dir ${DATA_DIR} \
  --initial-advertise-peer-urls https://${IP}:2380 \
  --listen-peer-urls https://${IP}:2380 \
  --listen-client-urls https://${IP}:2379,https://127.0.0.1:2379 \
  --advertise-client-urls https://${IP}:2379 \
  --initial-cluster-token ${INITIAL_CLUSTER_TOKEN} \
  --initial-cluster ${INITIAL_CLUSTER} \
  --initial-cluster-state ${INITIAL_CLUSTER_STATE} \
  --peer-auto-tls \
  --client-cert-auth \
  --trusted-ca-file=${SSL_DIR}/etcd.ca \
  --cert-file=${SSL_DIR}/server.cert \
  --key-file=${SSL_DIR}/server.key "

echo "options: ${OPTIONS}" >> ${LOGFILE}

netstat -anop | grep -E ":2379 |:2380 " >> ${LOGFILE}
[ $? -eq 0 ] && { echo "2379 or 2380 in use now" | tee -a ${LOGFILE}; exit 1; }

[ ! -e "${sourceFilePath}/etcd" ] && { echo 'etcd not found' | tee -a ${LOGFILE}; exit 1; }
chmod +x ${sourceFilePath}/etcd
${sourceFilePath}/etcd ${OPTIONS} >> ${LOGFILE} 2>&1 &

ps -ef | grep "${sourceFilePath}/etcd --name ${NAME}" | grep -v grep >> ${LOGFILE}

echo "Bootstrap etcd cluster finish" | tee -a ${LOGFILE}

#curl --cacert /etc/ssl/etcd/etcd.ca --cert /etc/ssl/etcd/client.cert --key /etc/ssl/etcd/client.key -L https://127.0.0.1:2379/v2/keys/foo -XPUT -d value=bar -k
#curl --cacert /etc/ssl/etcd/etcd.ca --cert /etc/ssl/etcd/client.cert --key /etc/ssl/etcd/client.key -L https://127.0.0.1:2379/v2/keys/foo -k

