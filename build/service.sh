#!/bin/bash


currentDir=$(pwd)


sudo tee /etc/systemd/system/cloud-manager.service > /dev/null << EOF
[Unit]
Description=Cloud SDN Management Service
After=ovn-controller.service
Requires=ovn-controller.service

[Service]
Type=simple
WorkingDirectory=$currentDir
ExecStart=$currentDir/cloud-manager
Restart=on-failure
RestartSec=5s
User=root

[Install]
WantedBy=multi-user.target
EOF


sudo systemctl daemon-reload
sudo systemctl enable --now cloud-manager
