#!/usr/bin/env bash

# the idea hear is to generate prebuilt binaries of kilcord to work on
# the killcord website. This will expand over time (I see rasberry pi 
# in the near future)
#
# the format for adding this to 
# killcord.io/killcord/0.0.1-alpha/killcord_0.0.1-alpha_macos_64bit.zip

# exit on error
set -e

# setup boiler plate
PROJECT="killcord"
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

case "$(uname -m)" in
	x86_64) CODE_ARCH="amd64" ;;
	x86) CODE_ARCH="386" ;;
	i686) CODE_ARCH="386" ;;
	i386) CODE_ARCH="386" ;;
	aarch64) CODE_ARCH="arm64" ;;
	armv5*) CODE_ARCH="arm5" ;;
	armv6*) CODE_ARCH="arm6" ;;
	armv7*) CODE_ARCH="arm7" ;;
	*) echo "unsupported architecture for this tool, exiting" && exit 1
esac

# go to the cmd directory
cd $PROJECT

# build it
echo "building killcord for $CODE_PLATFORM $CODE_ARCH"
go build

# grab the version
PROJECT_VERSION="$(./killcord version)"

# name it
ZIP_NAME=$PROJECT\_$PROJECT_VERSION\_$CODE_PLATFORM\_$CODE_ARCH.zip

# zip it
zip $ZIP_NAME $PROJECT

# upload it
aws s3 cp $ZIP_NAME $BUCKET/$PROJECT/$PROJECT_VERSION/$ZIP_NAME --acl public-read

# clean it
rm $PROJECT $ZIP_NAME
