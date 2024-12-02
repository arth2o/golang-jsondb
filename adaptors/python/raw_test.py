#!/usr/bin/env python3

from json_storage_connection import JsonStorageConnection
from config.config_manager import ConfigManager
import json
import time

def run_basic_tests(storage):
    print("\nRunning Basic Functionality Tests:")
    
    # Test PING
    response = storage.get("PING")
    print(f"PING test: {'PASSED' if response == 'PONG' else 'FAILED'}")
    
    # Test basic SET/GET
    storage.set("test:basic:1", "Hello World")
    response = storage.get("test:basic:1")
    print(f"Basic SET/GET test: {'PASSED' if response == 'Hello World' else 'FAILED'}")
    
    # Test DELETE
    deleted = storage.delete("test:basic:1")
    response = storage.get("test:basic:1")
    print(f"DELETE test: {'PASSED' if response is None else 'FAILED'}")
    
    # Test non-existent key
    response = storage.get("test:nonexistent")
    print(f"Non-existent key test: {'PASSED' if response is None else 'FAILED'}")

def run_data_type_tests(storage):
    print("\nRunning Data Type Tests:")
    
    test_cases = {
        'test:string': {
            'value': "Simple string test",
            'type': str.__name__
        },
        'test:integer': {
            'value': 42,
            'type': int.__name__
        },
        'test:float': {
            'value': 3.14159,
            'type': float.__name__
        },
        'test:boolean': {
            'value': True,
            'type': bool.__name__
        },
        'test:null': {
            'value': None,
            'type': type(None).__name__
        },
        'test:special_chars': {
            'value': "Special chars: !@#$%^&*()",
            'type': str.__name__
        }
    }
    
    for key, test in test_cases.items():
        storage.set(key, test['value'])
        retrieved = storage.get(key)
        
        type_match = isinstance(retrieved, type(test['value']))
        value_match = retrieved == test['value']
        
        print(f"{key}: {'PASSED' if type_match and value_match else 'FAILED'} "
              f"(Type: {type(retrieved).__name__}, Expected: {test['type']}) - "
              f"{str(test['value'])}")
        
        storage.delete(key)

def run_ttl_tests(storage):
    print("\nRunning TTL and Expiration Tests:")
    
    # Test TTL setting
    storage.set("test:ttl:1", "Expires in 2 seconds", 2)
    ttl = storage.ttl("test:ttl:1")
    print(f"TTL test (should be ~2): {'PASSED' if 0 < ttl <= 2 else 'FAILED'} (TTL: {ttl})")
    
    # Test expiration
    print("Waiting for key to expire...")
    time.sleep(3)
    expired = storage.get("test:ttl:1")
    print(f"Expiration test: {'PASSED' if expired is None else 'FAILED'}")
    
    # Test no expiration
    storage.set("test:ttl:2", "No expiration", -1)
    ttl = storage.ttl("test:ttl:2")
    print(f"No expiration test: {'PASSED' if ttl == -1 else 'FAILED'}")
    storage.delete("test:ttl:2")

def run_complex_data_tests(storage):
    print("\nRunning Complex Data Structure Tests:")
    
    test_cases = {
        'test:array:simple': {
            'data': ['apple', 'banana', 'orange'],
            'desc': 'Simple array'
        },
        'test:array:assoc': {
            'data': {'name': 'John', 'age': 30, 'city': 'New York'},
            'desc': 'Associative array'
        },
        'test:nested:deep': {
            'data': {
                'user': {
                    'profile': {
                        'name': 'Jane Doe',
                        'settings': {
                            'theme': 'dark',
                            'notifications': True,
                            'preferences': {
                                'language': 'en',
                                'timezone': 'UTC'
                            }
                        }
                    }
                }
            },
            'desc': 'Deeply nested structure'
        },
        'test:mixed:types': {
            'data': {
                'string': 'text',
                'number': 42,
                'float': 3.14,
                'boolean': True,
                'null': None,
                'array': [1, 2, 3],
                'object': {'key': 'value'}
            },
            'desc': 'Mixed data types'
        }
    }
    
    for key, test in test_cases.items():
        print(f"\nTesting {test['desc']}:")
        
        storage.set(key, test['data'])
        retrieved = storage.get(key)
        
        matches = json.dumps(test['data'], sort_keys=True) == json.dumps(retrieved, sort_keys=True)
        print(f"Data integrity: {'PASSED' if matches else 'FAILED'}")
        
        if not matches:
            print(f"Original: {json.dumps(test['data'])}")
            print(f"Retrieved: {json.dumps(retrieved)}")
        
        storage.delete(key)

def main():
    print("Starting JsonStorage connection test...\n")
    
    # Debug: Show all config values
    print("Loading configuration...")
    all_config = ConfigManager.all()
    print(f"Available configuration keys: {', '.join(all_config.keys())}")
    
    try:
        # Create and connect
        storage = JsonStorageConnection('localhost', 5555)
        storage.connect()
        
        # Run all test suites
        run_basic_tests(storage)
        run_data_type_tests(storage)
        run_ttl_tests(storage)
        run_complex_data_tests(storage)
        
    except Exception as e:
        print(f"\nError: {str(e)}")
    finally:
        if 'storage' in locals():
            storage.close()

if __name__ == "__main__":
    main()