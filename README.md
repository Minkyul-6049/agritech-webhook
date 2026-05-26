# Agritech Monitoring Webhook Server

A specialized monitoring middleware built with Go, designed to bridge Grafana alerts and automated notification systems for smart farm infrastructure.

## Key Features

* **Real-time Alert Processing:** Handles incoming HTTP POST webhooks from Grafana.
* **Dynamic Messaging:** Parses alert payloads to deliver context-aware notifications via Telegram.
* **Automated Response:** Triggers cooling or irrigation system logic based on farm environmental thresholds.
* **Secure Configuration:** Decoupled credentials management using runtime environment variables (Anti-Hardcoding Pattern).

## System Architecture

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
