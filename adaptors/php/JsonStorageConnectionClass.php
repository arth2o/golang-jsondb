<?php

require_once __DIR__ . '/config/ConfigManager.php';

/**
 * JsonStorageConnection Class
 * 
 * Provides a TCP socket-based connection to the JsonStorage server for storing
 * and retrieving data with optional TTL support.
 * 
 * @package JsonStorage
 * @author Your Name
 * @version 1.0.0
 */
class JsonStorageConnection {
    /** @var resource|null TCP socket resource */
    private $socket = null;
    
    /** @var string Hostname of the JsonStorage server */
    private $host;
    
    /** @var int Port number of the JsonStorage server */
    private $port;
    
    /** @var bool Connection status flag */
    private $connected = false;
    
    /** @var bool Authentication status flag */
    private $authenticated = false;
    
    /** @var bool Enable/disable debug output */
    private $debug = true;
    
    /** @var int Socket timeout in seconds */
    private $timeout = 5;

    /**
     * Constructor initializes the connection parameters
     * 
     * @param string $host Server hostname (default: 'localhost')
     * @param int|null $port Server port (default: from config or 5555)
     */
    public function __construct(string $host = 'localhost', ?int $port = null) {
        $this->host = $host;
        $this->port = $port ?? (int)ConfigManager::get('PORT', 5555);
        
        if ($this->debug) {
            echo "Initializing connection to {$this->host}:{$this->port}\n";
        }
    }

    /**
     * Handles server authentication using configured password
     * 
     * @throws Exception If authentication fails
     * @return bool True if authentication successful
     */
    private function handleAuthentication(): bool {
        if ($this->authenticated) {
            return true;
        }

        // Read initial server response
        $response = $this->read();
        if (trim($response) !== "AUTH_REQUIRED") {
            $this->authenticated = true;
            return true;
        }

        // Get password from ConfigManager and strip quotes
        $password = trim(ConfigManager::get('SERVER_PASSWORD'), "'\"");
        if (!$password) {
            throw new Exception("Authentication required but SERVER_PASSWORD not configured");
        }

        if ($this->debug) {
            echo "Attempting authentication...\n";
        }

        // Send AUTH command with password
        $authCommand = "AUTH " . $password;
        $this->write($authCommand);
        $response = $this->read();

        if (strpos($response, "OK") === false) {
            throw new Exception("Authentication failed: " . trim($response));
        }

        $this->authenticated = true;
        return true;
    }

    /**
     * Verifies authentication status before command execution
     * 
     * @throws Exception If not authenticated
     */
    private function verifyAuthentication(): void {
        if (!$this->authenticated) {
            throw new Exception("Not authenticated with server");
        }
    }

    /**
     * Writes data to the socket with error handling
     * 
     * @param string $data Data to write to the socket
     * @throws Exception If write operation fails
     * @return void
     */
    private function write(string $data): void {
        $data .= "\n";
        $written = socket_write($this->socket, $data, strlen($data));
        if ($written === false) {
            throw new Exception("Failed to write to socket: " . socket_strerror(socket_last_error()));
        }
    }

    /**
     * Reads data from the socket with error handling
     * 
     * @throws Exception If read operation fails
     * @return string Response from server
     */
    private function read(): string {
        $response = socket_read($this->socket, 1024);
        if ($response === false) {
            throw new Exception("Failed to read from socket: " . socket_strerror(socket_last_error()));
        }
        return trim($response);
    }

    /**
     * Sets a value in the storage with optional TTL
     * 
     * @param string $key Key to store the value under
     * @param mixed $value Value to store
     * @param int $ttl Time-to-live in seconds (-1 for no expiration)
     * @return bool Success status
     */
    public function set(string $key, $value, int $ttl = -1): bool {
        if (!$this->connected) {
            throw new Exception("Not connected");
        }
        $this->verifyAuthentication();

        // Handle different data types
        if (is_bool($value)) {
            $value = $value ? "true" : "false";
        } elseif (is_null($value)) {
            $value = "null";
        } elseif (is_array($value) || is_object($value)) {
            $value = json_encode($value);
        } elseif (is_numeric($value)) {
            // Don't quote numbers
            $value = (string)$value;
        } else {
            // Quote strings
            $value = '"' . addslashes($value) . '"';
        }

        // Build SET command
        $command = "SET {$key} {$value}";
        
        if ($this->debug) {
            echo "Sending command: " . $command . "\n";
        }
        
        $this->write($command);
        $response = $this->read();
        
        if (trim($response) !== 'OK') {
            return false;
        }

        // Handle TTL if specified
        if ($ttl > 0) {
            $this->write("EXPIRE {$key} {$ttl}");
            $expireResponse = $this->read();
            return trim($expireResponse) === 'OK';
        }

        return true;
    }

    /**
     * Gets a value from storage
     * 
     * @param string $key Key to retrieve
     * @return mixed Retrieved value
     */
    public function get(string $key) {
        if (!$this->connected) {
            throw new Exception("Not connected");
        }
        $this->verifyAuthentication();

        if ($key === "PING") {
            $this->write("PING");
            return $this->read();
        }

        $this->write("GET {$key}");
        $response = $this->read();

        if ($response === 'nil') {
            return null;
        }

        // Remove surrounding quotes if present
        $response = trim($response, '"');
        
        // Handle special values
        if ($response === 'null') {
            return null;
        } elseif ($response === 'true') {
            return true;
        } elseif ($response === 'false') {
            return false;
        }
        
        // Try to decode JSON if it looks like JSON
        if (($response[0] ?? '') === '{' || ($response[0] ?? '') === '[') {
            $decoded = json_decode($response, true);
            if (json_last_error() === JSON_ERROR_NONE) {
                return $decoded;
            }
        }

        // Handle numeric values
        if (is_numeric($response)) {
            return strpos($response, '.') !== false ? (float)$response : (int)$response;
        }

        return $response;
    }

    /**
     * Gets the TTL (time-to-live) for a key
     * 
     * @param string $key Key to check
     * @return int TTL in seconds, -1 if no expiration, -2 if key doesn't exist
     */
    public function ttl(string $key): int {
        if (!$this->connected) {
            throw new Exception("Not connected");
        }
        $this->verifyAuthentication();

        $this->write("TTL {$key}");
        $response = $this->read();
        
        if ($response === 'nil') {
            return -2;  // Key doesn't exist
        }
        
        return (int)$response;
    }

    /**
     * Establishes connection to the JsonStorage server and handles authentication
     * 
     * @throws Exception If connection or authentication fails
     * @return bool Success status
     */
    public function connect(): bool {
        if ($this->debug) {
            echo "Creating socket...\n";
        }
        
        // Create TCP/IP socket
        $this->socket = @socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        if ($this->socket === false) {
            $error = socket_last_error();
            throw new Exception("Failed to create socket: " . socket_strerror($error));
        }

        // Set socket timeout options
        socket_set_option($this->socket, SOL_SOCKET, SO_RCVTIMEO, array('sec' => $this->timeout, 'usec' => 0));
        socket_set_option($this->socket, SOL_SOCKET, SO_SNDTIMEO, array('sec' => $this->timeout, 'usec' => 0));

        if ($this->debug) {
            echo "Socket created successfully. Attempting to connect...\n";
        }

        // Attempt connection
        $result = @socket_connect($this->socket, $this->host, $this->port);
        if ($result === false) {
            $error = socket_last_error();
            throw new Exception("Failed to connect: " . socket_strerror($error));
        }

        $this->connected = true;

        // Handle authentication immediately after connection
        try {
            $this->handleAuthentication();
        } catch (Exception $e) {
            $this->close();
            throw $e;
        }

        if ($this->debug) {
            echo "Connection and authentication successful.\n";
        }

        return true;
    }

    /**
     * Closes the connection and cleans up resources
     * 
     * @return void
     */
    public function close(): void {
        if ($this->socket && $this->connected) {
            socket_close($this->socket);
            $this->connected = false;
            $this->socket = null;
        }
    }

    /**
     * Destructor ensures connection is properly closed
     */
    public function __destruct() {
        $this->close();
    }

    /**
     * Deletes a key from storage
     * 
     * @param string $key Key to delete
     * @throws Exception If not connected
     * @return bool True if key was deleted, false if key didn't exist
     */
    public function del(string $key): bool {
        if (!$this->connected) {
            throw new Exception("Not connected");
        }
        $this->verifyAuthentication();

        $this->write("DEL {$key}");
        $response = $this->read();

        // Server returns "1" for successful deletion, "0" for key not found
        return $response === "1";
    }

    /**
     * Adds debug bytes helper method
     * 
     * @param string $str String to debug
     * @return string Debugged string
     */
    private function debugBytes(string $str): string {
        $bytes = [];
        for ($i = 0; $i < strlen($str); $i++) {
            $bytes[] = sprintf('0x%02X[%s]', ord($str[$i]), $str[$i]);
        }
        return implode(' ', $bytes);
    }
}