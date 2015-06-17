#!/bin/bash

gox
for os in darwin freebsd linux netbsd; do
	for arch in amd64 386; do
		mkdir -p /tmp/$1/$os/$arch
		mv ./s3cp_${os}_${arch} /tmp/$1/$os/$arch/s3cp
	done
done

mkdir -p /tmp/$1/linux/arm
mv ./s3cp_linux_arm /tmp/$1/linux/arm/s3cp
mkdir -p /tmp/$1/freebsd/arm
mv ./s3cp_freebsd_arm /tmp/$1/freebsd/arm/s3cp
mkdir -p /tmp/$1/windows/386
mv ./s3cp_windows_386.exe /tmp/$1/windows/386/s3cp.exe
mkdir -p /tmp/$1/windows/amd64
mv ./s3cp_windows_amd64.exe /tmp/$1/windows/amd64/s3cp.exe

mv /tmp/$1/darwin /tmp/$1/osx

zip 
