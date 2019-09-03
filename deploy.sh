#! /bin/sh
# evel "$(ssh-agent -s)"
# ssh-add ~/.ssh/id_rsa
git pull http://github.com/ndcinfra/eventreward.git
go get -u github.com/ndcinfra/eventreward
go build
