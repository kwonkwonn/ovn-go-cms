#!/bin/bash


sudo tee /etc/systemd/system/cloud-manager.service > /dev/null << EOF
[Unit]
Description=Cloud SDN Management Service
After=ovn-controller.service
Requires=ovn-controller.service

[Service]
Type=simple
ExecStart=/home/kws/Desktop/cloud-manager/cloud-manager
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF


sudo systemctl daemon-reload
sudo systemctl enable --now ovn-go
