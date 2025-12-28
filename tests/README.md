# Order Placement Load Testing for Ha-Soranu

Focused load testing for the Ha-Soranu food delivery platform's order placement flow.

## Overview

This test evaluates the performance and scalability of the end-to-end order placement path:
- **Client** → **API Gateway** → **Restaurant Service** → **PostgreSQL** → **Kafka**.

## Prerequisites

### 1. Environment Setup

Ensure all services are running:

```bash
# Start services with Tilt (from project root)
cd /home/tamirat-dejene/Documents/dis-sys/ha-soranu
tilt up
```

Verify services are healthy:
```bash
curl http://localhost:8080/health
```

### 2. Python Environment

Create and activate virtual environment:

```bash
cd tests
make setup
source venv/bin/activate
```

## Quick Start

### 1. Seed Test Data

Before running tests, populate the database with test users and restaurants:

```bash
make seed
```

### 2. Run Order Load Test

```bash
# With Locust Web UI
make run

# Or headless via script
./run_tests.sh all
```

## Metrics and Results

Reports are generated in the `reports/` directory.
- `order_placement_report.html`: Visual summary of latencies, throughput, and error rates.

## Configuration

- `config/config.yaml`: Adjust user counts, spawn rates, and target latencies.
- `load/locustfile_order_placement.py`: Customize the user behavior and task weights.
