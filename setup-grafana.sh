cd
git clone https://github.com/hyperledger/sawtooth-core.git
cd sawtooth-core/docker
docker build . -f grafana/sawtooth-stats-grafana -t sawtooth-stats-grafana

cd
curl -sL https://repos.influxdata.com/influxdb.key |  sudo apt-key add -
sudo apt-add-repository "deb https://repos.influxdata.com/ubuntu xenial stable"
sudo apt-get update
sudo apt-get install telegraf
sudo cp voting/telegraf.conf /etc/telegraf/telegraf.conf
sudo systemctl restart telegraf