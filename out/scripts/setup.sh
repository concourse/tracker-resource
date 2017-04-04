#!/bin/bash

set -e

DIR=$1
ACTUAL_STORY_ID=$2
ACTUAL_DELIVERED_STORY_ID=$3

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
	git commit -m "add file [#223456 Finishes]" --allow-empty
	git commit -m "add file [Finished #323456]" --allow-empty
	git commit -m "add file [Finish #423456, #789456]" --allow-empty
	git commit -m "add file [Completes #523456]" --allow-empty
	git commit -m "add file [Completed #623456]" --allow-empty
	git commit -m "add file [Complete #723456]" --allow-empty
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
	git commit -m "fix bug [#223457 fixes]" --allow-empty
	git commit -m "fix bug [fixed #323457]" --allow-empty
	git commit -m "fix bug [fix #423457]" --allow-empty
	git commit -m "make bad change that gets rejected later [Complete #444444]" --allow-empty
	git commit -m "fix previously rejected story [Complete #555555]" --allow-empty
	git commit -m "fix previously rejected story badly so it gets rejected again [Complete #666666]" --allow-empty
popd

if [ ! -z ${ACTUAL_STORY_ID} ] || [ ! -z ${ACTUAL_DELIVERED_STORY_ID} ] ; then
	mkdir -p $DIR/middle/git3
	pushd $DIR/middle/git3
		git init

		git config user.email "concourse@example.com"
		git config user.name "Concourse Tracker Resource"
	popd
fi

if [ ! -z ${ACTUAL_STORY_ID} ]; then
	echo "ACTUAL_STORY_ID: ${ACTUAL_STORY_ID}"
	# git3: git with a vengence directory
	mkdir -p $DIR/middle/git3
	pushd $DIR/middle/git3
		echo "bugfix" > file.txt
		git add file.txt
		git commit -m "fix bug

		[fixes #${ACTUAL_STORY_ID}]"
	popd
fi

if [ ! -z ${ACTUAL_DELIVERED_STORY_ID} ]; then
	echo "ACTUAL_DELIVERED_STORY_ID: ${ACTUAL_DELIVERED_STORY_ID}"
	# git3: git with a vengence directory
	mkdir -p $DIR/middle/git3
	pushd $DIR/middle/git3
		echo "feature" > file.txt
		git add file.txt
		git commit -m "add feature

		[finishes #${ACTUAL_DELIVERED_STORY_ID}]"
	popd
fi
