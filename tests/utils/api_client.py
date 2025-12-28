import requests # type: ignore
import logging
from typing import Dict, Any

logger = logging.getLogger(__name__)

class APIClient:
    """Simplest possible API client for Ha-Soranu tests."""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()

    def post(self, endpoint: str, json_data: Dict) -> Dict:
        """Generic POST helper."""
        response = self.session.post(f"{self.base_url}{endpoint}", json=json_data, timeout=30)
        response.raise_for_status()
        return response.json()

    def health_check(self) -> bool:
        """Check if API Gateway is healthy."""
        try:
            return self.session.get(f"{self.base_url}/health", timeout=5).status_code == 200
        except:
            return False

    def register_restaurant(self, data: Dict) -> Dict:
        """Register a restaurant."""
        return self.post('/api/v1/restaurants/register', data)

    def list_restaurants(self, lat: float = 9.0, lon: float = 38.0, radius: float = 100.0) -> Dict:
        """List restaurants."""
        response = self.session.post(
            f"{self.base_url}/api/v1/restaurants",
            params={"Latitude": lat, "Longitude": lon, "RadiusKm": radius},
            timeout=10
        )
        response.raise_for_status()
        return response.json()
