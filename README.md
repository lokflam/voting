# Voting
Blockchain voting system prototype

## Development environment (optional)
Docker & Docker Compose
```bash
# docker
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce
sudo usermod -a -G docker ${USER}

# docker-compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.23.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# login again
```

## Keys
The keys used as example in this prototype:  
Public key: `02caf7596b9b8d9c27da437df19f6de57e94a5dad903ec0cd15549179377164379`  
Private key: `793b849da13c4829693fa555c54686e44951f227637929e0997bd1b67705ecec`

## Example network
Host 1: Node 1  
Host 2: Node 2  
Host 3: Node 3  
Host 4: Node 4  
Host 5: Load balancer

## Setup
For every host running as a node in blockchain network
1. `cd voting`
1. Set public IP: `export HOST=$(curl ifconfig.me)`
1. Set peers' public IP, for example:  
    Host 1: 1.1.1.1  
    Host 2: 2.2.2.2  
    Host 3: 3.3.3.3  
    Host 4: 4.4.4.4  
    
    Command for host 1: `export PEERS=tcp://2.2.2.2:8800,tcp://3.3.3.3:8800,tcp://4.4.4.4:8800`  
    Command for host 2: `export PEERS=tcp://1.1.1.1:8800,tcp://3.3.3.3:8800,tcp://4.4.4.4:8800`  
    Command for host 3: `export PEERS=tcp://1.1.1.1:8800,tcp://2.2.2.2:8800,tcp://4.4.4.4:8800`  
    Command for host 4: `export PEERS=tcp://1.1.1.1:8800,tcp://2.2.2.2:8800,tcp://3.3.3.3:8800`
1. Run Docker
    1. For first node: `docker-compose -f sawtooth-voting-0.yaml up -d`
    1. For other nodes: `docker-compose -f sawtooth-voting-1.yaml up -d`
    1. Shut down
        1. `docker-compose -f sawtooth-voting-0.yaml down -v`
        1. `docker-compose -f sawtooth-voting-1.yaml down -v`
    
## Set permission
Set permission for transaction processor (set using the first node)
1. `docker cp sawtooth-validator-default:/etc/sawtooth/keys/validator.priv .`
1. `docker cp validator.priv sawtooth-shell-default:.`
1. `docker exec -it sawtooth-shell-default bash`
1. `sawset proposal create --key validator.priv sawtooth.identity.allowed_keys=$(cat ~/.sawtooth/keys/root.pub) --url http://rest-api:8008`
1. ```bash
   sawtooth identity policy create vo_policy \
     --url http://rest-api:8008 \
     "PERMIT_KEY 02caf7596b9b8d9c27da437df19f6de57e94a5dad903ec0cd15549179377164379" \
     "DENY_KEY *"
   ```
1. ```bash
   sawtooth identity role create \
     --url http://rest-api:8008 \
     transactor.transaction_signer.voting-organizer vo_policy
   ```

## Setup load balancer
Only need to setup on one host
1. Install nginx
    ```bash
    sudo apt-get update
    sudo apt-get install nginx
    ```
1. Add config for load balancing to `/etc/nginx/conf.d/`
    1. Can use `/nginx.conf` in this project as example

## Load test with wrk
1. Install wrk (referencing from https://github.com/wg/wrk/wiki/Installing-wrk-on-Linux)
    ```bash
    sudo apt-get install build-essential libssl-dev git -y
    git clone https://github.com/wg/wrk.git wrk
    cd wrk
    make
    # move the executable to somewhere in your PATH, ex:
    sudo cp wrk /usr/local/bin
    ```
1. Create a new vote
1. Add ballot to vote using script
    ```bash
    # host = 1.1.1.1:9009
    # private key = 793b849da13c4829693fa555c54686e44951f227637929e0997bd1b67705ecec
    # vote id = 1
    # number of ballots = 100
    # this command adds ballots with ballot code = "1", "2", ..., "100"
    . addballot.sh 1.1.1.1:9009 793b849da13c4829693fa555c54686e44951f227637929e0997bd1b67705ecec 1 100
    ```
1. Run load test
    ```bash
    # run load test on vote id = 1 & number of ballots = 100
    env vote_id=1 num=100 wrk -t1 -c10 -s wrk.lua http://1.1.1.1:9009/ballot/cast
    ```
    
