FROM debian:wheezy

# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN useradd --create-home --home /home/mongodb --uid 8000 --user-group --shell /bin/bash mongodb

# wheezy is no longer provided by deb.debian.org, replaced with archive.debian.org
RUN sed -i 's|^deb http://deb.debian.org/debian|deb http://archive.debian.org/debian|' /etc/apt/sources.list
# remove wheezy-updates (no longer exist, as there are no planned point releases)
RUN sed -i '/wheezy-updates/d' /etc/apt/sources.list
RUN sed -i '/deb http:\/\/security.debian.org\/debian-security wheezy\/updates main/d' /etc/apt/sources.list

RUN apt-get -o Acquire::Check-Valid-Until=false update

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
		ca-certificates \
		numactl \
	&& rm -rf /var/lib/apt/lists/*

# grab gosu for easy step-down from root (https://github.com/tianon/gosu/releases)
ENV GOSU_VERSION 1.10
# grab "js-yaml" for parsing mongod's YAML config files (https://github.com/nodeca/js-yaml/releases)
ENV JSYAML_VERSION 3.10.0

RUN set -ex; \
	\
	apt-get update; \
	apt-get install -y --no-install-recommends \
		wget \
	; \
	rm -rf /var/lib/apt/lists/*; \
	\
	dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
	wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch"; \
	wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc"; \
	export GNUPGHOME="$(mktemp -d)"; \
	gpg --keyserver keyserver.ubuntu.com --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \
	gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \
	command -v gpgconf && gpgconf --kill all || :; \
	rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc; \
	chmod +x /usr/local/bin/gosu; \
	gosu nobody true; \
	\
	wget -O /js-yaml.js "https://github.com/nodeca/js-yaml/raw/${JSYAML_VERSION}/dist/js-yaml.js"; \
# TODO some sort of download verification here
	\
	apt-get purge -y --auto-remove wget

RUN mkdir /var/run/mongodb

RUN install -d -m 0755 -o mongodb -g mongodb /var/run/mongodb

ENV GPG_KEYS BE0EB6AA24F3C978
# pub   4096R/AAB2461C 2014-02-25 [expires: 2016-02-25]
#       Key fingerprint = DFFA 3DCF 326E 302C 4787  673A 01C4 E7FA AAB2 461C
# uid                  MongoDB 2.6 Release Signing Key <packaging@mongodb.com>
# pub   4096R/EA312927 2015-10-09 [expires: 2017-10-08]
#       Key fingerprint = 42F3 E95A 2C4F 0827 9C49  60AD D68F A50F EA31 2927
# uid                  MongoDB 3.2 Release Signing Key <packaging@mongodb.com>
# https://docs.mongodb.com/manual/tutorial/verify-mongodb-packages/#download-then-import-the-key-file
RUN set -ex; \
	export GNUPGHOME="$(mktemp -d)"; \
	for key in $GPG_KEYS; do \
		gpg --keyserver keyserver.ubuntu.com --recv-keys "$key"; \
	done; \
	gpg --export $GPG_KEYS > /etc/apt/trusted.gpg.d/mongodb.gpg; \
	command -v gpgconf && gpgconf --kill all || :; \
	rm -r "$GNUPGHOME"; \
	apt-key list

# Allow build-time overrides (eg. to build image with MongoDB Enterprise version)
# Options for MONGO_PACKAGE: mongodb-org OR mongodb-enterprise
# Options for MONGO_REPO: repo.mongodb.org OR repo.mongodb.com
# Example: docker build --build-arg MONGO_PACKAGE=mongodb-enterprise --build-arg MONGO_REPO=repo.mongodb.com .
ARG MONGO_PACKAGE=mongodb-org
ARG MONGO_REPO=repo.mongodb.org
ENV MONGO_PACKAGE=${MONGO_PACKAGE} MONGO_REPO=${MONGO_REPO}

ENV MONGO_MAJOR 3.0
ENV MONGO_VERSION 3.0.15

RUN echo "deb http://$MONGO_REPO/apt/debian wheezy/${MONGO_PACKAGE}/$MONGO_MAJOR main" | tee "/etc/apt/sources.list.d/${MONGO_PACKAGE%-unstable}.list"

RUN set -x \
	&& apt-get update \
	&& apt-get install -y --force-yes \
		${MONGO_PACKAGE}=$MONGO_VERSION \
		${MONGO_PACKAGE}-server=$MONGO_VERSION \
		${MONGO_PACKAGE}-shell=$MONGO_VERSION \
		${MONGO_PACKAGE}-mongos=$MONGO_VERSION \
		${MONGO_PACKAGE}-tools=$MONGO_VERSION \
	&& rm -rf /var/lib/apt/lists/* \
	&& rm -rf /var/lib/mongodb \
	&& mv /etc/mongod.conf /etc/mongod.conf.orig

COPY config/etc/mongod.conf /etc/mongod.conf

RUN mkdir -p /var/lib/mongo  \
	&& chown -R mongodb:mongodb /var/lib/mongo

VOLUME /var/lib/mongo

COPY config/etc/mongodb-keyfile /var/lib/

RUN chown -R mongodb:mongodb /var/lib/mongo

RUN chown -R mongodb:mongodb /var/lib/mongodb-keyfile

RUN chown -R mongodb:mongodb /etc/mongod.conf

RUN chmod 400 /var/lib/mongodb-keyfile

COPY entrypoint.sh /entrypoint.sh

RUN chown -R mongodb:mongodb /etc/mongod.conf

RUN chmod 777 /entrypoint.sh

RUN chown -R mongodb:mongodb /entrypoint.sh

RUN chmod o+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

USER mongodb

EXPOSE 27017

CMD ["mongod"]