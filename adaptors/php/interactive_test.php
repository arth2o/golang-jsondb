<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

echo "Starting interactive test client...\n";

// Debug: Show all config values
echo "Loading configuration...\n";
$allConfig = ConfigManager::all();
echo "Available configuration keys: " . implode(", ", array_keys($allConfig)) . "\n";

try {
    $storage = new JsonStorageConnection('localhost', 5555);
    $storage->connect();
    
    echo "Connected to server. Type commands (quit to exit):\n";
    echo "Available commands: SET key value, GET key, DELETE key, TTL key, PING\n";
    
    while (true) {
        echo "> ";
        $command = trim(fgets(STDIN));
        
        if ($command === 'quit') {
            break;
        }
        
        try {
            $parts = explode(' ', $command);
            $cmd = strtoupper($parts[0]);
            
            switch ($cmd) {
                case 'SET':
                    if (count($parts) < 3) {
                        echo "Error: SET requires key and value parameters\n";
                        continue;
                    }
                    $key = $parts[1];
                    $value = implode(' ', array_slice($parts, 2));
                    $result = $storage->set($key, $value);
                    echo "Result: " . ($result ? "OK" : "FAILED") . "\n";
                    break;
                    
                case 'GET':
                    if (count($parts) < 2) {
                        echo "Error: GET requires a key parameter\n";
                        continue;
                    }
                    $result = $storage->get($parts[1]);
                    echo "Result: " . ($result === null ? "nil" : $result) . "\n";
                    break;
                    
                case 'DELETE':
                    if (count($parts) < 2) {
                        echo "Error: DELETE requires a key parameter\n";
                        continue;
                    }
                    $result = $storage->delete($parts[1]);
                    echo "Result: " . ($result ? "OK" : "FAILED") . "\n";
                    break;
                    
                case 'TTL':
                    if (count($parts) < 2) {
                        echo "Error: TTL requires a key parameter\n";
                        continue;
                    }
                    $result = $storage->ttl($parts[1]);
                    echo "Result: " . $result . "\n";
                    break;
                    
                case 'PING':
                    $result = $storage->get('PING');
                    echo "Result: " . $result . "\n";
                    break;
                    
                default:
                    echo "Unknown command. Available commands: SET, GET, DELETE, TTL, PING\n";
            }
        } catch (Exception $e) {
            echo "Command error: " . $e->getMessage() . "\n";
            continue;
        }
    }
    
} catch (Exception $e) {
    echo "\nERROR: " . $e->getMessage() . "\n";
    echo "Stack trace:\n" . $e->getTraceAsString() . "\n";
} finally {
    if (isset($storage)) {
        echo "\nClosing connection...\n";
        $storage->close();
    }
}

echo "Connection closed.\n";