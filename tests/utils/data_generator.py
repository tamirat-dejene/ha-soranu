import random
import uuid
from typing import Dict, List
from faker import Faker # type: ignore

class DataGenerator:
    """Generates data exactly matching Ha-Soranu API requirements."""
    
    def __init__(self, seed: int = None):
        self.fake = Faker()
        if seed is not None:
            random.seed(seed)
            Faker.seed(seed)

    def generate_restaurant(self) -> Dict:
        """Matches RegisterRestaurantDTO."""
        return {
            "name": f"{self.fake.company()} Kitchen",
            "email": self.fake.email(),
            "secret_key": str(uuid.uuid4()),
            "latitude": round(random.uniform(9.0, 9.1), 6),
            "longitude": round(random.uniform(38.7, 38.8), 6),
            "menus": [
                {
                    "name": self.fake.word().capitalize() + " Special",
                    "description": self.fake.sentence(),
                    "price": float(round(random.uniform(50, 500), 2))
                }
                for _ in range(random.randint(5, 10))
            ]
        }

    def generate_order(self, customer_id: str, restaurant_id: str, menu_items: List[Dict]) -> Dict:
        """Matches PlaceOrderDTO."""
        if not menu_items:
            raise ValueError(f"Cannot generate order for restaurant {restaurant_id} with no menu items")
            
        num_items = random.randint(1, min(3, len(menu_items)))
        selected = random.sample(menu_items, num_items)
        
        return {
            "customer_id": customer_id,
            "restaurant_id": restaurant_id,
            "items": [
                {
                    "item_id": item["item_id"],
                    "quantity": random.randint(1, 5)
                }
                for item in selected
            ]
        }

    def generate_search_area(self) -> Dict:
        """Matches ListRestaurantsDTO."""
        return {
            "latitude": 9.0192,
            "longitude": 38.7525,
            "radius_km": 10.0
        }
