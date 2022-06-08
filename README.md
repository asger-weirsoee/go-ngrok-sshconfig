# go-ngrok-sshconfig
So Ngrok only allows for one active machine per user for a free account.
But I want to use ngrok for multiple machines.

Introducing `go-ngrok-sshconfig`.

Now, by creating multiple accounts and then adding the api keys for each account in this .token file.

`go-ngrok-sshconfig` will then create add these services to you .ssh/config file, allowing you to easy connect to your machines.


## Format of .tokens file

The .token file is a .csv file, where the first index should be your ngrok api key, and the second index should be the name of the machine that will also be appended to the .ssh/config file as the host

Every line in the .token file is a new set of api key and machine name.

```
token,token_name,agc
token2,token_name2,cba

```

would result in a hosts file looking like:


```
Host token_name
    HostName (ngrok url from api key)
    Port (ngrok port from api key)
    User abc
    ServerAliveInterval 300
    ServerAliveCountMax 3

Host token_name2
    HostName (ngrok url from api key)
    Port (ngrok port from api key)
    User cba
    ServerAliveInterval 300
    ServerAliveCountMax 3
```


## Examples of server setup:

### /etc/systemd/system/ngrok.service:
```
[Unit]
Description=Ngrok
After=network.service

[Service]
type=simple
User=maskine
WorkinDirectory=/home/maskine
ExecStart=/usr/bin/ngrok start --all --config="/home/maskine/ngrok_config/config.yml"
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### ~/ngrok_config/config.yml:
```
authtoken: authtoken
tunnels:
    default:
        proto: tcp
        addr: 22
version: "2"
region: eu
```
