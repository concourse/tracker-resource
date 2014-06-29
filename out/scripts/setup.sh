#!/bin/bash

set -e

DIR=$1

# random file
pushd $DIR
	touch random-file
popd

# empty directory
pushd $DIR
	mkdir empty
popd

# non-git directory
mkdir -p $DIR/notgit

pushd $DIR/notgit
	echo "filez" > file.txt
popd

# git directory
mkdir -p $DIR/git
pushd $DIR/git
	git init

	touch file.txt
	git add file.txt
	git commit -m "add file [finishes #123456]"

	echo "bugfix" > file.txt
	git add file.txt
	git commit -m "fix bug [fixes #123457]"
popd