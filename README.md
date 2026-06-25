# 🚜 Agritech Edge Observability Pipeline

![Pipeline Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Kubernetes](https://img.shields.io/badge/K3s-Edge_Cluster-326CE5?logo=kubernetes)
![Grafana](https://img.shields.io/badge/Grafana-v10-F46800?logo=grafana)

A lightweight, resilient telemetry data pipeline designed for resource-constrained agricultural edge environments (Farm-nodes). This project demonstrates the migration from a legacy prototyping tool to a production-grade monitoring stack.

## 📸 Dashboard Overview
<img width="1920" height="1080" alt="Screenshot 2026-06-23 234514" src="https://github.com/user-attachments/assets/6d9c21e3-dd10-4cf0-860c-dfcb7557a539" />


## 🏗️ Architecture Flow
The architecture is designed to ensure zero data loss and minimal resource overhead on the edge node:

1. **Data Ingestion (Golang Webhook):** A highly optimized custom webhook written in Go receives sensor data (Temperature, Humidity, Soil Moisture).
2. **Resource Management:** Runs as a standalone `systemd` daemon on the Linux edge node, ensuring high availability with `Restart=on-failure` auto-healing mechanisms. 
3. **Metrics Storage (Prometheus):** Scrapes and stores time-series data efficiently within the K3s embedded cluster.
4. **Visualization & Alerting (Grafana):** Provides real-time, actionable insights with clear thresholds. Critical alerts (e.g., node downtime, abnormal temperatures) are routed directly to **Telegram** for immediate operational response.

## 🚀 Key Achievements
* **Extreme Resource Efficiency:** The Go ingestion webhook is optimized to consume less than **3 MiB of Memory** and **0.15% CPU**, proving its suitability for low-powered edge devices.
* **Production-Grade Reliability:** Transitioned from Node-RED to a robust Grafana/Prometheus stack, implementing OS-level process management (systemd) to prevent alerting failures (e.g., resolving 404 Webhook delivery errors).
* **Actionable Insights:** Conf
### System Topology
```mermaid
graph TD
    subgraph Farm-Node [Farm-Node: 192.168.202.131]
        SD[Golang Sensor Daemon] -->|Auto-Trigger| LCL[Local Control Logic]
        LCL -->|Local Actuation| SV[Sprinkler / Ventilation]
    end

    subgraph Monitor-Node [Monitor-Node: 192.168.202.132]
        Prom[Prometheus Server] -->|2. Query Metrics| GD[Grafana Dashboard]
        Prom -->|3. Query Metrics| NR[Node-RED HMI]
        GD -->|4. Trigger Alert| GWS[Golang Webhook Server]
        NR -->|6. Manual Override| SV
    end

    SD -->|1. HTTP Pull| Prom
    GWS -->|5. Format and Push| TG[Telegram App]

    style Farm-Node fill:#f9f,stroke:#333,stroke-width:2px
    style Monitor-Node fill:#bbf,stroke:#333,stroke-width:2px
    style TG fill:#85C1E9,stroke:#333,stroke-width:1px
