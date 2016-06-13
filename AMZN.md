##### Installing on Amazon Linux.

SSH to instance.

```
ssh -i "AMI.pem" ec2-user@54.172.11.229
```

Install golang.

```
sudo yum -y --enablerepo=epel install golang-bin
```

Ensure everything is working okay.

```
mkdir ~/.go
export GOPATH=~/.go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH
go get golang.org/x/tour/gotour
gotour
```