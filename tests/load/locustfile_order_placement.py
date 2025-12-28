import sys
import uuid
import random
import yaml # type: ignore
from pathlib import Path
from locust import HttpUser, task, between # type: ignore

sys.path.insert(0, str(Path(__file__).parent.parent))
from utils.data_generator import DataGenerator

# Load seeded data for quick access
data_path = Path(__file__).parent.parent / 'config' / 'seeded_data.yaml'
if data_path.exists():
    with open(data_path, 'r') as f:
        seeded_data = yaml.safe_load(f) or {"restaurants": []}
else:
    seeded_data = {"restaurants": []}

class OrderPlacementUser(HttpUser):
    wait_time = between(1, 2)

    def on_start(self):
        self.gen = DataGenerator()
        self.customer_id = "46b0302c-b7fb-4ee9-ba67-58cbe868be4c"
        self.restaurants = seeded_data.get("restaurants", [])

    @task
    def place_order(self):
        if not self.restaurants: return
        
        target = random.choice(self.restaurants)
        order = self.gen.generate_order(
            customer_id=self.customer_id,
            restaurant_id=target["restaurant_id"],
            menu_items=target["menus"]
        )

        with self.client.post("/api/v1/restaurants/orders", json=order, catch_response=True) as resp:
            if resp.status_code in [200, 201]:
                resp.success()
            else:
                resp.failure(f"Order failed: {resp.status_code} - {resp.text}")