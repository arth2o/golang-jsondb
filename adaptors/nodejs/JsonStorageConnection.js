const net = require('net');
const { ConfigManager } = require('./config/ConfigManager');

class JsonStorageConnection {
    constructor(host = 'localhost', port = null) {
        this.host = host;
        this.port = port ?? ConfigManager.get('PORT', 5555);
        this.socket = null;
        this.connected = false;
        this.authenticated = false;
        this.debug = true;
        this.timeout = 5000; // milliseconds
        this.maxListeners = 20;
    }

    async connect() {
        if (this.debug) {
            console.log(`Creating socket connection to ${this.host}:${this.port}...`);
        }

        return new Promise((resolve, reject) => {
            this.socket = new net.Socket();
            
            this.socket.setTimeout(this.timeout);

            this.socket.setMaxListeners(this.maxListeners);

            this.socket.on('timeout', () => {
                reject(new Error('Socket timeout'));
                this.close();
            });

            this.socket.connect(this.port, this.host, async () => {
                this.connected = true;
                try {
                    await this.handleAuthentication();
                    resolve(true);
                } catch (error) {
                    this.close();
                    reject(error);
                }
            });

            this.socket.on('error', (error) => {
                reject(error);
                this.close();
            });
        });
    }

    async handleAuthentication() {
        if (this.authenticated) return true;

        const response = await this.read();
        if (response.trim() !== 'AUTH_REQUIRED') {
            this.authenticated = true;
            return true;
        }

        const password = ConfigManager.get('SERVER_PASSWORD')?.replace(/['"]/g, '');
        if (!password) {
            throw new Error('Authentication required but SERVER_PASSWORD not configured');
        }

        if (this.debug) {
            console.log('Attempting authentication...');
        }

        await this.write(`AUTH ${password}`);
        const authResponse = await this.read();

        if (!authResponse.includes('OK')) {
            throw new Error(`Authentication failed: ${authResponse.trim()}`);
        }

        this.authenticated = true;
        return true;
    }

    verifyAuthentication() {
        if (!this.authenticated) {
            throw new Error('Not authenticated with server');
        }
    }

    write(data) {
        return new Promise((resolve, reject) => {
            if (!this.socket) {
                reject(new Error('Socket not connected'));
                return;
            }

            this.socket.write(data + '\n', (error) => {
                if (error) reject(error);
                else resolve();
            });
        });
    }

    read() {
        return new Promise((resolve, reject) => {
            if (!this.socket) {
                reject(new Error('Socket not connected'));
                return;
            }

            const dataHandler = (data) => {
                cleanup();
                resolve(data.toString().trim());
            };

            const errorHandler = (error) => {
                cleanup();
                reject(error);
            };

            const cleanup = () => {
                this.socket.removeListener('data', dataHandler);
                this.socket.removeListener('error', errorHandler);
            };

            this.socket.once('data', dataHandler);
            this.socket.once('error', errorHandler);
        });
    }

    async set(key, value, ttl = -1) {
        if (!this.connected) throw new Error('Not connected');
        this.verifyAuthentication();

        // Handle different data types
        let processedValue;
        if (typeof value === 'string') {
            // For strings, escape quotes and wrap in quotes
            processedValue = `"${value.replace(/"/g, '\\"')}"`;
        } else if (typeof value === 'boolean') {
            processedValue = value ? 'true' : 'false';
        } else if (value === null) {
            processedValue = 'null';
        } else if (typeof value === 'object') {
            processedValue = JSON.stringify(value);
        } else if (typeof value === 'number') {
            processedValue = value.toString();
        }

        const command = `SET ${key} ${processedValue}`;
        
        if (this.debug) {
            console.log(`Sending command: ${command}`);
            console.log('Processed value:', processedValue);
        }

        await this.write(command);
        const response = await this.read();

        if (response.trim() !== 'OK') {
            return false;
        }

        if (ttl > 0) {
            await this.write(`EXPIRE ${key} ${ttl}`);
            const expireResponse = await this.read();
            return expireResponse.trim() === 'OK';
        }

        return true;
    }

    async get(key) {
        if (!this.connected) throw new Error('Not connected');
        this.verifyAuthentication();

        if (key === 'PING') {
            await this.write('PING');
            return await this.read();
        }

        await this.write(`GET ${key}`);
        const response = await this.read();

        if (response === 'nil') return null;

        // Remove surrounding quotes if present
        if (response.startsWith('"') && response.endsWith('"')) {
            return response.slice(1, -1).replace(/\\"/g, '"');
        }

        // Handle other types
        if (response === 'true') return true;
        if (response === 'false') return false;
        if (response === 'null') return null;
        
        // Try parsing as number
        if (!isNaN(response)) {
            return response.includes('.') ? parseFloat(response) : parseInt(response);
        }
        
        // Try parsing as JSON for objects and arrays
        try {
            if ((response.startsWith('{') && response.endsWith('}')) ||
                (response.startsWith('[') && response.endsWith(']'))) {
                return JSON.parse(response);
            }
        } catch (e) {
            // If JSON parsing fails, return as is
        }

        return response;
    }

    async ttl(key) {
        if (!this.connected) throw new Error('Not connected');
        this.verifyAuthentication();

        await this.write(`TTL ${key}`);
        const response = await this.read();

        if (response === 'nil') return -2;
        return parseInt(response);
    }

    async del(key) {
        if (!this.connected) throw new Error('Not connected');
        this.verifyAuthentication();

        await this.write(`DEL ${key}`);
        const response = await this.read();

        return response === '1';
    }

    close() {
        if (this.socket && this.connected) {
            this.socket.destroy();
            this.connected = false;
            this.socket = null;
        }
    }
}

module.exports = { JsonStorageConnection };