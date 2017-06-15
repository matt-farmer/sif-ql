#!/bin/bash

set -e

CWD=`pwd`

# code is required for build, but no longer built as separate server
# is launched from within the naprr application
# echo "Downloading nats-streaming-server"
# go get github.com/nats-io/nats-streaming-server

do_build() {
	mkdir -p $OUTPUT
	cd $CWD
	cd ./app/qlserver
	go get
	GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUTPUT/$HARNESS
	cd ..
	rsync -a qlserver/naplan_schema.graphql qlserver/public $OUTPUT/
}

# do_shells() {
# 	cd $CWD
# 	cp bin/gonias.sh $OUTPUT/
# 	cp bin/stopnias.sh $OUTPUT/
# }

# do_bats() {
# 	cd $CWD
# 	cp bin/gonias.bat $OUTPUT/
# 	cp bin/stopnias.bat $OUTPUT/
# }

# do_upx() {
# 	upx $OUTPUT/$GNATS
# 	upx $OUTPUT/$HARNESS
# }

# do_goupx() {
# 	goupx $OUTPUT/$GNATS
# 	goupx $OUTPUT/$HARNESS
# }

do_zip() {
	cd $OUTPUT
	cd ..
	zip -qr ../$ZIP qlserver
	cd $CWD
}

build_mac64() {
	# MAC OS X (64 only)
	echo "Building Mac binaries..."
	GOOS=darwin
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Mac/qlserver
	# GNATS=nats-streaming-server
	HARNESS=qlserver
	ZIP=qlserver-Mac.zip
	do_build
	#do_upx
	# do_shells
	do_zip
	echo "...all Mac binaries built..."
}


build_windows64() {
	# WINDOWS 64
	echo "Building Windows64 binaries..."
	GOOS=windows
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win64/qlserver
	# GNATS=nats-streaming-server.exe
	HARNESS=qlserver.exe
	ZIP=qlserver-Win64.zip
	do_build
	#do_upx
	# do_bats
	do_zip
	echo "...all Windows64 binaries built..."
}

build_windows32() {
	# WINDOWS 32
	echo "Building Windows32 binaries..."
	GOOS=windows
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Win32/qlserver
	# GNATS=nats-streaming-server.exe
	HARNESS=qlserver.exe
	ZIP=qlserver-Win32.zip
	do_build
	#do_upx
	# do_bats
	do_zip
	echo "...all Windows32 binaries built..."
}

build_linux64() {
	# LINUX 64
	echo "Building Linux64 binaries..."
	GOOS=linux
	GOARCH=amd64
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux64/qlserver
	# GNATS=nats-streaming-server
	HARNESS=qlserver
	ZIP=qlserver-Linux64.zip
	do_build
	#do_goupx
	# do_shells
	do_zip
	echo "...all Linux64 binaries built..."
}

build_linux32() {
	# LINUX 32
	echo "Building Linux32 binaries..."
	GOOS=linux
	GOARCH=386
	LDFLAGS="-s -w"
	OUTPUT=$CWD/build/Linux32/qlserver
	# GNATS=nats-streaming-server
	HARNESS=qlserver
	ZIP=qlserver-Linux32.zip
	do_build
	#do_goupx
	# do_shells
	do_zip
	echo "...all Linux32 binaries built..."
}

# TODO ARM
# GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator

build_mac64
build_windows64
#build_windows32
#build_linux64
#build_linux32

