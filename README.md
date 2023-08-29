# remote-sign-svr
Support for cryptographic signature services for different chains (EVM, TRON, DOGE, CKB).

# Prerequisites

* Ubuntu 18.04 or newer (2C4G)
* MYSQL >= 8.0
* GO version >= 1.17.10

## Install & Run

### Source Compile

```bash
# pull the code
git clone https://github.com/dotbitHQ/remote-sign-svr.git

# modify the configuration file (Database Configuration)
cd remote-sign-svr
cp config/config.example.yaml config/config.yaml 
vim config/config.yaml 

# compile and run the service
make svr
./remote_sign_svr --config=config/config.yaml

# compile and run the management client
make cli
./remote_sign_cli
```
