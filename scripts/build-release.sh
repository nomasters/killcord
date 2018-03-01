#!/usr/bin/env bash

set -e

# the idea hear is to generate prebuilt binaries of kilcord to work on
# the killcord website. This will expand over time (I see rasberry pi 
# in the near future)
#
# the format for adding this to 
# killcord.io/downloads/killcord/VERSION/PLATFORM/ARCHITECTURE/killcord

PROJ="killcord"
BUCKET="s3://killcord.io"
CODE_PLATFORM="null"
CODE_ARCH="null"

# get platform name
case "$(uname -s)" in
    Linux*)     
		CODE_PLATFORM="linux"
		;;
    Darwin*)    
		CODE_PLATFORM="macos"
    	;;
    *)
		echo "unsupported OS, exiting"
		exit 1
esac

# get supported architectures
case "$(uname -m)" in
    x86_64)     
		CODE_ARCH="64-bit"
		;;
    *)
		echo "unsupported architecture for this tool, exiting"
		exit 1
esac

cd $PROJ
echo "building killcord for $CODE_PLATFORM $CODE_ARCH"
go build

KILLCORD_VERSION="$(./killcord version)"

aws s3 cp $PROJ $BUCKET/downloads/killcord/$KILLCORD_VERSION/$CODE_PLATFORM/$CODE_ARCH/$PROJ --acl public-read
rm killcord
