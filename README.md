# 🚜 Agritech Edge Observability Pipeline

![Pipeline Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Kubernetes](https://img.shields.io/badge/K3s-Edge_Cluster-326CE5?logo=kubernetes)
![Grafana](https://img.shields.io/badge/Grafana-v10-F46800?logo=grafana)

A lightweight, resilient telemetry data pipeline designed for resource-constrained agricultural edge environments (Farm-nodes). This project demonstrates the migration from a legacy prototyping tool to a production-grade monitoring stack.

## 📸 Dashboard Overview
<img width="1920" height="1080" alt="Screenshot 2026-06-23 234514" src="https://github.com/user-attachments/assets/a483951d-2e76-4853-bae4-d522e484b2c3" />

## 🏗️ Architecture Flow
The architecture is designed to ensure zero data loss and minimal resource overhead on the edge node. 

### System Topology
```mermaid
graph LR
    subgraph "Farm Node (192.168.202.131)"
        Sensors["🌱 IoT Sensors<br>(Temp, Humidity, Soil)"]
    end

    subgraph "Monitor Node (192.168.202.132)"
        subgraph "K3s Embedded Cluster"
            DB[("🗄️ Time-Series DB<br>(InfluxDB/Prometheus)")]
            Grafana["📊 Grafana v10<br>(Dashboard & Alerting)"]
        end
        
        subgraph "Linux OS (systemd)"
            GoWebhook{"⚙️ Go Webhook Receiver<br>(Port: 8080/webhook)"}
        end
    end

    Telegram["📱 Telegram API<br>(Mobile Alerts)"]

    Sensors -- "Raw Data" --> DB
    DB -- "Query" --> Grafana
    Grafana -- "Alert Trigger<br>(HTTP POST)" --> GoWebhook
    GoWebhook -- "Format & Route" --> Telegram

    %% Styling
    classDef k3s fill:#326ce5,stroke:#fff,stroke-width:2px,color:#fff;
    classDef os fill:#333,stroke:#fff,stroke-width:2px,color:#fff;
    classDef external fill:#0088cc,stroke:#fff,stroke-width:2px,color:#fff;
    classDef farm fill:#4caf50,stroke:#fff,stroke-width:2px,color:#fff;

    class Grafana,DB k3s;
    class GoWebhook os;
    class Telegram external;
    class Sensors farm;
