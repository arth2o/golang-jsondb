import os
from typing import Dict, Optional

class EnvLoader:
    @staticmethod
    def load(environment: str) -> Dict[str, str]:
        """
        Loads environment variables from a .env file
        
        Args:
            environment (str): Development, production, or testing
            
        Returns:
            Dict[str, str]: Loaded environment variables
            
        Raises:
            FileNotFoundError: If the .env file cannot be found
        """
        env_file = os.path.join(
            os.path.dirname(os.path.dirname(os.path.dirname(__file__))),
            'jsondb',
            f'.env.{environment}'
        )
        
        if not os.path.exists(env_file):
            raise FileNotFoundError(f"Environment file not found: {env_file}")
            
        config = {}
        with open(env_file, 'r') as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith('#'):
                    continue
                    
                key, value = line.split('=', 1)
                key = key.strip()
                value = value.strip().strip("'").strip('"')
                config[key] = value
                
        return config
