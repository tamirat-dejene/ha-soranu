#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m'

# Minimal Setup
if [ ! -d "venv" ]; then
    python3 -m venv venv
fi
source venv/bin/activate
pip install -q -r requirements.txt
mkdir -p reports

case "$1" in
    seed)
        python3 utils/seed_data.py
        ;;
    run)
        locust -f load/locustfile_order_placement.py --host http://localhost:8080 --users 20 --spawn-rate 2 --run-time 1m --headless
        ;;
    all)
        python3 utils/seed_data.py
        locust -f load/locustfile_order_placement.py --host http://localhost:8080 --users 20 --spawn-rate 2 --run-time 1m --headless
        ;;
    *)
        echo "Usage: $0 {seed|run|all}"
        exit 1
        ;;
esac
