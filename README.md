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

### Run by Supervisor

```bash
# install supervisor first
sudo apt install supervisor

# pull the code
git clone https://github.com/dotbitHQ/remote-sign-svr.git

# modify the configuration file (Database Configuration)
cd remote-sign-svr
cp config/config.example.yaml config/config.yaml 
vim config/config.yaml 
mkdir logs

# compile and run the service
make svr

# updating the configuration of supervisor
vim /etc/supervisor/conf.d/remote_sign_svr.conf
supervisorctl update
supervisorctl status

# compile and run the management client
make cli
./remote_sign_cli

```

#### /etc/supervisor/conf.d/remote_sign_svr.conf
```yaml
[program:remote_sign_svr]
# Project Catalog  
directory = /mnt/server/remote-sign-svr
command = /mnt/server/remote-sign-svr/remote_sign_svr --config=/mnt/server/remote-sign-svr/config/config.yaml

autostart=true                ; start at supervisord start (default: true)
autorestart=true
user=root                   ; setuid to this UNIX account to run the program
startsecs=2
startretries=3

redirect_stderr=true          ; redirect proc stderr to stdout (default false)
stdout_logfile=/mnt/server/remote-sign-svr/logs/out.log        ; stdout log path, NONE for none; default AUTO
stdout_logfile_maxbytes=100MB   ; max # logfile bytes b4 rotation (default 50MB)
stdout_logfile_backups=20     ; # of stdout logfile backups (default 10)
stdout_capture_maxbytes=100MB   ; number of bytes in 'capturemode' (default 0)
stdout_events_enabled=false   ; emit events on stdout writes (default false)
stderr_logfile=/mnt/server/remote-sign-svr/logs/err.log        ; stderr log path, NONE for none; default AUTO
stderr_logfile_maxbytes=100MB   ; max # logfile bytes b4 rotation (default 50MB)
stderr_logfile_backups=20     ; # of stderr logfile backups (default 10)
stderr_capture_maxbytes=100MB   ; number of bytes in 'capturemode' (default 0)
stderr_events_enabled=false   ; emit events on stderr writes (default false)
```
