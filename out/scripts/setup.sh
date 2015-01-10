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

	git config user.email "concourse@example.com"
  	git config user.name "Concourse Tracker Resource"

	touch file.txt
	git add file.txt
	git commit -m "add file [Finishes #123456]"
popd

# git2: git harder directory
mkdir -p $DIR/middle/git2
pushd $DIR/middle/git2
	git init

	git config user.email "concourse@example.com"
  	git config user.name "Concourse Tracker Resource"

	echo "bugfix" > file.txt
	git add file.txt
	git commit -m "fix bug

	[fixes #123457]"
popd
