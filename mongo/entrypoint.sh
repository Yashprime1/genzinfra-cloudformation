#!/bin/bash
set -Eeuo pipefail

if [ "${1:0:1}" = '-' ]; then
	set -- mongod "$@"
fi

originalArgOne="$1"

CONFIGFILE="/etc/mongod.conf"
DAEMON_OPTS=" --config $CONFIGFILE"
SYSCONFIG="/etc/sysconfig/mongod"
# allow the container to be started with `--user`
# all mongo* commands should be dropped to the correct user

# you should use numactl to start your mongod instances, including the config servers, mongos instances, and any clients.
# https://docs.mongodb.com/manual/administration/production-notes/#configuring-numa-on-linux

DAEMON_OPTS=${DAEMON_OPTS:-"--config $CONF"}


PIDFILEPATH=$(grep pidFilePath "$CONFIGFILE" | awk -F ':' {'print $2'} | awk -F '#' {'print $1'})

mongod=${MONGOD-/usr/bin/mongod}

MONGO_USER=mongodb
MONGO_GROUP=mongodb

if [ -f "$SYSCONFIG" ]; then
    . "$SYSCONFIG"
fi

echo $PIDFILEPATH

PIDDIR=`dirname $PIDFILEPATH`

id

if [ ! -d $PIDDIR ]; then
install -d -m 0755 -o $MONGO_USER -g $MONGO_GROUP $PIDDIR
fi

# Make sure the pidfile does not exist
if [ -f $PIDFILEPATH ]; then
  echo "Error starting mongod. $PIDFILEPATH exists."
  RETVAL=1
  return
fi

# Recommended ulimit values for mongod or mongos
# See http://docs.mongodb.org/manual/reference/ulimit/#recommended-settings

echo -n $"Starting mongod: "

exec numactl --interleave=all mongod -f /etc/mongod.conf


