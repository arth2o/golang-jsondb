<?php

require_once __DIR__ . '/EnvLoader.php';

/**
 * Class ConfigManager
 * 
 * Manages configuration settings for the application by loading and caching
 * environment-specific variables.
 * 
 * @package JsonDB\Config
 */
class ConfigManager {
    /** @var array|null Cached configuration values */
    private static ?array $config = null;
    
    /** @var string Current environment (development, production, testing) */
    private static string $environment = 'development';
    
    /**
     * Sets the current environment and resets the configuration cache.
     * 
     * @param string $env The environment to set (development, production, testing)
     * @return void
     */
    public static function setEnvironment(string $env): void {
        self::$environment = $env;
        self::$config = null; // Reset config to force reload on next access
    }
    
    /**
     * Retrieves a specific configuration value by key.
     * 
     * If the configuration hasn't been loaded yet, it will load it from the
     * environment file. If the key doesn't exist, returns the default value.
     * 
     * @param string $key The configuration key to retrieve
     * @param mixed $default The default value if the key doesn't exist
     * @return mixed The configuration value or default if not found
     */
    public static function get(string $key, mixed $default = null): mixed {
        // Load configuration if not already loaded
        if (self::$config === null) {
            try {
                self::$config = EnvLoader::load(self::$environment);
            } catch (Exception $e) {
                error_log("Failed to load environment: " . $e->getMessage());
                return $default;
            }
        }
        
        return self::$config[$key] ?? $default;
    }

    /**
     * Retrieves all configuration values.
     * 
     * If the configuration hasn't been loaded yet, it will load it from the
     * environment file. Returns an empty array if loading fails.
     * 
     * @return array All configuration values
     */
    public static function all(): array {
        // Load configuration if not already loaded
        if (self::$config === null) {
            try {
                self::$config = EnvLoader::load(self::$environment);
            } catch (Exception $e) {
                error_log("Failed to load environment: " . $e->getMessage());
                return [];
            }
        }
        
        return self::$config;
    }
}