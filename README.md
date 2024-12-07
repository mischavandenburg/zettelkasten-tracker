# An overengineered tracker for my Zettelkasten

Running it as a systemd timer on Linux:

```bash

# /etc/systemd/system/zkcount.service
[Unit]
Description=ZK Count Script
 
[Service]
Type=oneshot
ExecStart=/home/mischa/Repos/github.com/mischavandenburg/dotfiles/scripts/zkcount
User=mischa

[Install]
WantedBy=multi-user.target

# /etc/systemd/system/zkcount.timer
[Unit]
Description=Run ZK Count Script hourly

[Timer]
OnBootSec=5min
OnUnitActiveSec=1h
Unit=zkcount.service

[Install]
WantedBy=timers.target

```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now zkcount.timer
``
