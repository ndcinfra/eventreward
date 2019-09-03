#! /bin/sh
evel "$(ssh-agent -s)"
ssh-add ~/.ssh/id_rsa
git pull
go get -u github.com/ndcinfra/eventreward
go build
