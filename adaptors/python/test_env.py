#!/usr/bin/env python3

import os
import sys

def find_env_file(environment='development'):
    """Find the environment file by walking up directories"""
    current_dir = os.path.dirname(os.path.abspath(__file__))
    
    while current_dir != '/':
        # Try to find jsondb directory
        jsondb_dir = os.path.join(current_dir, 'jsondb')
        if os.path.exists(jsondb_dir):
            env_file = os.path.join(jsondb_dir, f'.env.{environment}')
            if os.path.exists(env_file):
                return env_file
        current_dir = os.path.dirname(current_dir)
    return None

def main():
    env_file = find_env_file()
    if env_file:
        print(f"Found environment file: {env_file}")
        with open(env_file, 'r') as f:
            print("\nEnvironment contents:")
            for line in f:
                if line.strip() and not line.startswith('#'):
                    print(line.strip())
    else:
        print("Environment file not found!")
        sys.exit(1)

if __name__ == "__main__":
    main()