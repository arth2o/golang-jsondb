#!/usr/bin/env python3

from json_storage_connection import JsonStorageConnection
import time

print("Starting basic commands example...\n")

try:
    # Initialize connection
    storage = JsonStorageConnection('localhost', 5555)
    print("Connecting to server...")
    storage.connect()
    
    # Basic commands demonstration
    print("\nRunning basic commands:")
    
    # PING test
    print("\n=== PING Test ===")
    response = storage.get("PING")
    print(f"PING response: {response}")
    
    # SET command
    print("\n=== SET Command ===")
    key = "example:test:1"
    value = "Hello from basic commands!"
    success = storage.set(key, value)
    print(f"SET result: {'SUCCESS' if success else 'FAILED'}")
    
    # GET command
    print("\n=== GET Command ===")
    retrieved = storage.get(key)
    print(f"GET result: {retrieved}")
    
    # SET with TTL
    print("\n=== SET with TTL ===")
    ttl_key = "example:ttl:1"
    ttl_value = "This will expire in 5 seconds"
    storage.set(ttl_key, ttl_value, 5)
    print(f"Initial TTL value: {storage.get(ttl_key)}")
    print(f"TTL remaining: {storage.ttl(ttl_key)} seconds")
    
    # Wait and check TTL
    time.sleep(2)
    print(f"After 2 seconds - TTL remaining: {storage.ttl(ttl_key)} seconds")
    
    # DELETE command
    print("\n=== DELETE Command ===")
    storage.delete(key)
    check_deleted = storage.get(key)
    print(f"After DELETE - value exists?: {'Yes' if check_deleted is not None else 'No'}")

except Exception as e:
    print(f"\nError: {str(e)}")
    import traceback
    print("Stack trace:")
    print(traceback.format_exc())
finally:
    if 'storage' in locals():
        storage.close()
        print("\nConnection closed.")