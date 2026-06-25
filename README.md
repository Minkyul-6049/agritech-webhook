# 🚜 Agritech Edge Observability Pipeline

![Pipeline Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Kubernetes](https://img.shields.io/badge/K3s-Edge_Cluster-326CE5?logo=kubernetes)
![Grafana](https://img.shields.io/badge/Grafana-v10-F46800?logo=grafana)

A specialized monitoring middleware and resilient telemetry data pipeline built with Go. Designed for resource-constrained agricultural edge environments (Farm-nodes), this project bridges Grafana alerts and automated notification systems for smart farm infrastructure.

## 📸 Dashboard Overview
<img width="1920" height="1080" alt="Screenshot 2026-06-23 234514" src="https://github.com/user-attachments/assets/444bcc4b-b365-41c9-9715-564d4bf99c75" />


## 🏗️ Architecture Flow

1. Farm sensors send data to the embedded Time-Series DB (Prometheus/InfluxDB).
2. Grafana monitors thresholds (e.g., Temp > 30°C).
3. Grafana triggers a Webhook to this Go Server.
4. The Go Server processes the logic and notifies via Telegram, while enabling local actuation.

### System Topology

```mermaid
graph TD
    subgraph FarmNode [Farm-Node: 192.168.202.131]
        SD[Golang Sensor Daemon] -->|Auto-Trigger| LCL[Local Control Logic]
        LCL -->|Local Actuation| SV[Sprinkler / Ventilation]
    end

    subgraph MonitorNode [Monitor-Node: 192.168.202.132]
        Prom[Prometheus Server] -->|2. Query Metrics| GD[Grafana Dashboard]
        Prom -->|3. Query Metrics| NR[Node-RED HMI]
        GD -->|4. Trigger Alert| GWS[Golang Webhook Server]
        NR -->|6. Manual Override| SV
    end

    SD -->|1. HTTP Pull| Prom
    GWS -->|5. Format and Push| TG[Telegram App]

    style FarmNode fill:#f9f,stroke:#333,stroke-width:2px
    style MonitorNode fill:#bbf,stroke:#333,stroke-width:2px
    style TG fill:#85C1E9,stroke:#333,stroke-width:1px
