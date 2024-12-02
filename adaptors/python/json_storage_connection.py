import socket
import json
from typing import Any, Optional, Union
from config.config_manager import ConfigManager

class JsonStorageConnection:
    def __init__(self, host: str = 'localhost', port: Optional[int] = None):
        self.host = host
        self.port = port or int(ConfigManager.get('PORT', 5555))
        self.socket: Optional[socket.socket] = None
        self.connected = False
        self.authenticated = False
        self.debug = True
        self.timeout = 5
        
    def connect(self) -> bool:
        if self.debug:
            print("Creating socket...")
            
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.settimeout(self.timeout)
        
        try:
            self.socket.connect((self.host, self.port))
            self.connected = True
            self._handle_authentication()
            return True
        except Exception as e:
            raise ConnectionError(f"Failed to connect: {e}")
            
    def _handle_authentication(self) -> bool:
        if self.authenticated:
            return True
            
        response = self._read()
        if response.strip() != "AUTH_REQUIRED":
            self.authenticated = True
            return True
            
        password = ConfigManager.get('SERVER_PASSWORD', '').strip("'\"")
        if not password:
            raise Exception("Authentication required but SERVER_PASSWORD not configured")
            
        if self.debug:
            print("Attempting authentication...")
            
        self._write(f"AUTH {password}")
        response = self._read()
        
        if "OK" not in response:
            raise Exception(f"Authentication failed: {response.strip()}")
            
        self.authenticated = True
        return True
        
    def _write(self, data: str) -> None:
        if not self.socket:
            raise Exception("Not connected")
        self.socket.send(f"{data}\n".encode())
        
    def _read(self) -> str:
        if not self.socket:
            raise Exception("Not connected")
        return self.socket.recv(1024).decode()
        
    def set(self, key: str, value: Any, ttl: int = -1) -> bool:
        if not self.connected:
            raise Exception("Not connected")
            
        if isinstance(value, (dict, list)):
            value = json.dumps(value)
        
        command = f"SET {key} {value}"
        self._write(command)
        response = self._read()
        
        if response.strip() != 'OK':
            return False
            
        if ttl > 0:
            self._write(f"EXPIRE {key} {ttl}")
            expire_response = self._read()
            return expire_response.strip() == 'OK'
            
        return True
        
    def get(self, key: str) -> Any:
        if not self.connected:
            raise Exception("Not connected")
            
        if key == "PING":
            self._write("PING")
            return self._read().strip()
            
        self._write(f"GET {key}")
        response = self._read().strip()
        
        if response == 'nil':
            return None
            
        response = response.strip('"')
        
        if response == 'null':
            return None
        elif response == 'true':
            return True
        elif response == 'false':
            return False
            
        try:
            if response.startswith('{') or response.startswith('['):
                return json.loads(response)
        except:
            pass
            
        try:
            if '.' in response:
                return float(response)
            return int(response)
        except:
            pass
            
        return response
        
    def ttl(self, key: str) -> int:
        if not self.connected:
            raise Exception("Not connected")
            
        self._write(f"TTL {key}")
        response = self._read().strip()
        
        if response == 'nil':
            return -2
            
        return int(response)
        
    def delete(self, key: str) -> bool:
        if not self.connected:
            raise Exception("Not connected")
            
        self._write(f"DEL {key}")
        response = self._read().strip()
        return response == "1"
        
    def close(self) -> None:
        if self.socket and self.connected:
            self.socket.close()
            self.connected = False
            self.socket = None