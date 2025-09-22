# OVN-based Virtual Network Automation

A lightweight project to automate **virtual networking** for small-scale private cloud environments using **Open vSwitch (OVS)** and **OVN**.  
This project focuses on simplicity and efficiency, without relying on heavyweight SDN controllers.

---

## ✨ Features
- **Automated Port Management**  
  - Dynamically add OVS ports during VM creation.  

- **Isolated Subnet Allocation**  
  - Assign VXLAN tunnels for VMs requiring separate subnets.  

- **Gateway Integration**  
  - Use the **Control server’s OVS L3 switch** as a gateway for internal ↔ external connectivity.  

- **Lightweight Automation**  
  - Replace traditional controllers (e.g., OpenDaylight) with **gRPC APIs** + **Python bindings** for OVS.  

---

## 🏗️ Architecture Overview
```

+-------------------+  
| Client (VM API)   |  
+-------------------+  
|                   |  
| api               |  
v                   |  
+-------------------+ +--------------------+  
| Control Server |------| External Network |  
| (Automation) | GW | (Internet / LAN)     |  
| | +--------------------------------------+  
| - api gateway                            |  
| - OVS/OVN controller, northd Mgmt        |  
+---------+---------------------------------+  
|    control node    |  
| OVS Port / VXLAN   |  
v                    |  
+--------------------+  
| Compute Nodes      |  
| (VM Instances)     |  
+-------------------+
```


- **Client**: Sends VM creation/deletion requests.  
- **Control Server**:  
  - Runs OVS/OVN.  
  - Manages ports, VXLAN tunnels, and L3 gateway rules.  
  - Provides gRPC APIs for automation.  
- **Compute Nodes**: Host VMs connected to OVS bridges.  

---

## 🔧 Core Components

### 1. Control Server
- **OVS/OVN Integration**  
  - Handles port creation (`ovs-vsctl add-port`).  
  - Allocates VXLAN tunnels for network isolation.  
  - Exposes VM lifecycle hooks (create, delete).  
  - Communicates with compute nodes for sync.  

### 2. Network Layer
- **Geneve** for tenant-level segmentation.  


### 3. VM Management
- VM network interfaces dynamically attach to OVS bridges.  
- Each VM can be assigned to:  
  - **Default subnet** (shared)  
  - **Dedicated isolated ip subnet** 

---


## ⚠️ Notes
- Tested in debian based environment.  
- May require modification of filesystem paths and OVS configuration depending on distribution.  




🚀 A step toward **lightweight cloud networking automation** for developers and small private clouds.
