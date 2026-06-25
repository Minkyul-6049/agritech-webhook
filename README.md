
---

```markdown
# 🚜 Agritech Edge Observability Pipeline

![Pipeline Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Kubernetes](https://img.shields.io/badge/K3s-Edge_Cluster-326CE5?logo=kubernetes)
![Grafana](https://img.shields.io/badge/Grafana-v10-F46800?logo=grafana)

A specialized monitoring middleware and resilient telemetry data pipeline built with Go. Designed for resource-constrained agricultural edge environments (Farm-nodes), this project bridges Grafana alerts and automated notification systems for smart farm infrastructure.

## 📸 Dashboard Overview

<img width="1920" height="1080" alt="Screenshot 2026-06-23 234514" src="https://github.com/user-attachments/assets/ade8e93b-0691-4f85-890f-783198c25986" />


## 🏗️ Architecture Flow

1. Farm sensors send data to the embedded Time-Series DB (Prometheus/InfluxDB).
2. Grafana monitors thresholds (e.g., Temp > 30°C).
3. Grafana triggers a Webhook to this Go Server.
4. The Go Server processes the logic and notifies via Telegram, while enabling local actuation.

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
```

## 🛠️ Middleware Technical Specifications

### Key Features

* **Real-time Alert Processing:** Handles incoming HTTP POST webhooks from Grafana.
* **Dynamic Messaging:** Parses alert payloads to deliver context-aware notifications via Telegram.
* **Automated Response:** Triggers cooling or irrigation system logic based on farm environmental thresholds.
* **Secure Configuration:** Decoupled credentials management using runtime environment variables (Anti-Hardcoding Pattern).

### Tech Stack

* **Language:** Go 1.20+ (Optimized for low-latency & concurrency)
* **Monitoring/Alerting:** Grafana v10 & Prometheus
* **Notification API:** Telegram Bot API
* **Infrastructure:** Linux systemd (Daemonized with auto-recovery)

## 🚀 Key Achievements

* **Extreme Resource Efficiency:** The Go ingestion webhook is optimized to consume less than **3 MiB of Memory** and **0.15% CPU**, proving its suitability for low-powered edge devices.
* **Production-Grade Reliability:** Transitioned from Node-RED to a robust Grafana/Prometheus stack, implementing OS-level process management (systemd) to prevent alerting failures (e.g., resolving 404 Webhook delivery errors).
* **Actionable Insights:** Configured intuitive thresholds (Green/Orange/Red) to prevent alert fatigue and bridge the communication gap between IT infrastructure and agronomy operations.

## ⚙️ Setup Instructions

1. Configure `.env` with your `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`.
2. Build: `go build -o webhook-server main.go`
3. Run as a service using the provided `systemd` configuration (`Restart=on-failure` enabled).
