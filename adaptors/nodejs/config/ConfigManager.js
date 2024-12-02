const { EnvLoader } = require('./EnvLoader');

class ConfigManager {
    static config = null;
    static environment = process.env.NODE_ENV || 'development';

    static initialize() {
        if (!this.config) {
            try {
                this.config = EnvLoader.load(this.environment);
            } catch (error) {
                console.error('Failed to load configuration:', error.message);
                this.config = {};
            }
        }
    }

    /**
     * Get a configuration value by key
     * @param {string} key - Configuration key
     * @param {*} defaultValue - Default value if key not found
     * @returns {*} Configuration value or default value
     */
    static get(key, defaultValue = null) {
        if (!this.config) {
            this.initialize();
        }
        return this.config[key] ?? defaultValue;
    }

    /**
     * Get all configuration values
     * @returns {Object} All configuration values
     */
    static all() {
        if (!this.config) {
            this.initialize();
        }
        return { ...this.config };
    }

    /**
     * Set a configuration value
     * @param {string} key - Configuration key
     * @param {*} value - Configuration value
     */
    static set(key, value) {
        if (!this.config) {
            this.initialize();
        }
        this.config[key] = value;
    }

    /**
     * Check if a configuration key exists
     * @param {string} key - Configuration key
     * @returns {boolean} True if key exists
     */
    static has(key) {
        if (!this.config) {
            this.initialize();
        }
        return key in this.config;
    }
}

module.exports = { ConfigManager };