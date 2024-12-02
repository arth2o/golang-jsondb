import os
from typing import Dict, Any, Optional

class ConfigManager:
    _config: Optional[Dict[str, str]] = None
    _environment: str = 'development'
    
    @classmethod
    def find_env_file(cls) -> Optional[str]:
        current_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
        while current_dir != '/':
            jsondb_dir = os.path.join(current_dir, 'jsondb')
            if os.path.exists(jsondb_dir):
                env_file = os.path.join(jsondb_dir, f'.env.{cls._environment}')
                if os.path.exists(env_file):
                    return env_file
            current_dir = os.path.dirname(current_dir)
        return None

    @classmethod
    def load_config(cls) -> Dict[str, str]:
        env_file = cls.find_env_file()
        if not env_file:
            raise FileNotFoundError(f"Environment file not found for {cls._environment}")
            
        config = {}
        with open(env_file, 'r') as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith('#'):
                    continue
                key, value = line.split('=', 1)
                config[key.strip()] = value.strip().strip("'").strip('"')
        return config

    @classmethod
    def get(cls, key: str, default: Any = None) -> Any:
        if cls._config is None:
            try:
                cls._config = cls.load_config()
            except Exception as e:
                print(f"Failed to load environment: {e}")
                return default
        return cls._config.get(key, default)
    
    @classmethod
    def all(cls) -> Dict[str, str]:
        if cls._config is None:
            try:
                cls._config = cls.load_config()
            except Exception as e:
                print(f"Failed to load environment: {e}")
                return {}
        return cls._config
