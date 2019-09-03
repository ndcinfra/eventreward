#! /bin/sh
git pull http://github.com/ndcinfra/eventreward.git
go get -u github.com/ndcinfra/eventreward
go build
