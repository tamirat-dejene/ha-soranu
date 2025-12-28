import sys
import logging
import yaml # type: ignore
from pathlib import Path

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from utils.api_client import APIClient
from utils.data_generator import DataGenerator

logging.basicConfig(level=logging.INFO, format='%(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def seed_restaurants(count: int):
    client = APIClient()
    gen = DataGenerator()
    restaurants = []

    if not client.health_check():
        logger.error("API Gateway unreachable.")
        return

    logger.info("Checking for existing restaurants...")
    try:
        existing = client.list_restaurants()
        for r in existing.get("restaurants", []):
            if r.get("menus"):
                restaurants.append({
                    "restaurant_id": r["restaurant_id"],
                    "menus": r["menus"]
                })
        logger.info(f"Found {len(restaurants)} viable restaurants (with menus).")
    except Exception as e:
        logger.warning(f"Could not fetch existing restaurants: {e}")

    needed = count - len(restaurants)
    if needed > 0:
        logger.info(f"Seeding {needed} new restaurants...")
        for i in range(needed):
            try:
                data = gen.generate_restaurant()
                resp = client.register_restaurant(data)
                restaurants.append({
                    "restaurant_id": resp["restaurant_id"],
                    "menus": resp["menus"]
                })
                if (i+1) % 5 == 0: logger.info(f"Seeds: {i+1}/{needed}")
            except Exception as e:
                msg = str(e)
                if hasattr(e, 'response') and e.response is not None:
                    msg += f" - Body: {e.response.text}"
                logger.error(f"Failed seed: {msg}")
    
    if not restaurants:
        logger.error("No restaurants available to seed or use. Tests will fail.")
        return

    output_path = Path(__file__).parent.parent / 'config' / 'seeded_data.yaml'
    with open(output_path, 'w') as f:
        yaml.dump({"restaurants": restaurants}, f)
    logger.info(f"Seeding complete. Saved to {output_path}")

if __name__ == "__main__":
    seed_restaurants(25)
