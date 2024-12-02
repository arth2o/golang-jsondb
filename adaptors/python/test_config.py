#!/usr/bin/env python3

from config.config_manager import ConfigManager

def main():
    print("Testing configuration loading...")
    
    try:
        config = ConfigManager.all()
        print("\nLoaded configuration:")
        for key, value in config.items():
            print(f"{key}: {value}")
            
        # Test specific values
        port = ConfigManager.get('PORT')
        password = ConfigManager.get('SERVER_PASSWORD')
        
        print("\nSpecific values:")
        print(f"PORT: {port}")
        print(f"SERVER_PASSWORD: {password}")
        
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()