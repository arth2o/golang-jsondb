<?php

/**
 * Class EnvLoader
 * 
 * Handles loading and parsing of environment-specific configuration files.
 * Supports multiple environments (development, production, testing) and
 * parses .env files into key-value pairs.
 * 
 * @package JsonDB\Config
 */
class EnvLoader {
    /**
     * Loads and parses environment variables from a .env file.
     * 
     * Supports the following features:
     * - Skips empty lines and comments (lines starting with #)
     * - Handles quoted values (both single and double quotes)
     * - Strips surrounding quotes from values
     * - Throws exception if environment file is not found
     *
     * @param string $environment The environment to load (development, production, testing)
     * @return array Associative array of environment variables
     * @throws Exception If the .env file cannot be found or read
     */
    public static function load(string $environment): array {
        // Construct path to environment file
        $envFile = __DIR__ . '/../../../jsondb/.env.' . $environment;
        
        // Check if environment file exists
        if (!file_exists($envFile)) {
            throw new Exception("Environment file not found: $envFile");
        }
        
        $config = [];
        // Read file line by line, ignoring empty lines
        $lines = file($envFile, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
        
        foreach ($lines as $line) {
            // Skip comments
            if (strpos(trim($line), '#') === 0) {
                continue;
            }
            
            // Split line into key and value
            list($key, $value) = explode('=', $line, 2);
            $key = trim($key);
            $value = trim($value);
            
            // Remove surrounding quotes if present
            if ((substr($value, 0, 1) === "'" && substr($value, -1) === "'") ||
                (substr($value, 0, 1) === '"' && substr($value, -1) === '"')) {
                $value = substr($value, 1, -1);
            }
            
            $config[$key] = $value;
        }
        
        return $config;
    }
}