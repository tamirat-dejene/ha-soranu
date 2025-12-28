.PHONY: setup seed run all clean

VENV = venv
PYTHON = $(VENV)/bin/python3
PIP = $(VENV)/bin/pip

setup:
	python3 -m venv $(VENV)
	$(PIP) install -r requirements.txt
	mkdir -p reports

seed:
	$(PYTHON) utils/seed_data.py --users 100 --restaurants 50

run:
	$(VENV)/bin/locust -f load/locustfile_order_placement.py --host http://localhost:8080

all: seed run

clean:
	rm -rf $(VENV)
	rm -rf reports/*.html
	rm -rf reports/*.json
	rm -rf config/seeded_data.yaml
	find . -type d -name "__pycache__" -exec rm -rf {} +
