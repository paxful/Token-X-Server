Configuration
=============

**WARNING**

This server is NOT designed to be publicly exposed. It is a backend system and needs to be securely guarded because it contains private keys.

Caveat
---

I am not a sysops guy - any where you know better please correct

A parity node will require at least :

t2.medium, 50 GB SSD, 4 GB ram.

Install parity
--- 

**note : we decided to use QuikNode for now, later you may wish to control your own node**


`bash <(curl https://get.parity.io -kL)`

You can make the parity node separate from the Toke-X-Server

Install latest version of GO (was 1.8 at time of writing)
---

either 

```
sudo add-apt-repository ppa:gophers/archive
sudo apt update
sudo apt install golang-1.8 or later using 
```

or read the following

https://medium.com/@patdhlk/how-to-install-go-1-9-1-on-ubuntu-16-04-ee64c073cd79

Configure golang
---

add `/usr/lib/go-1.8/bin` to path
mkdir `~/go`
set `GOPATH=~/go`

Load GO-Ethereum 
---

clone `github.com/ethereum/go-ethereum` into `~/go/src/github.com/ethereum/go-ethereum`    

Load Libraries
---

```
go get github.com/spf13/viper
go get github.com/lib/pq
go get github.com/DaveAppleton/ether_go
go get github.com/DaveAppleton/etherUtils
go get github.com/DaveAppleton/etherdb
go get github.com/DaveAppleton/parityclient
```

Install Postgres
---

```
$ wget -q https://www.postgresql.org/media/keys/ACCC4CF8.asc -O - | sudo apt-key add -
$ sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
$ sudo apt-get update
$ sudo su - postgres <---- 
$ psql
```

Setup Databases
---

* when you get etherdb, there is a sql file included. This contains the sql to create the tables.
* This creates a database called etherdb
* config.json contains the connection parameters that should be used to connect to that db



