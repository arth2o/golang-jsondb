<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

echo "Starting basic connection test...\n";

// Debug: Show all config values
echo "Loading configuration...\n";
$allConfig = ConfigManager::all();
echo "Available configuration keys: " . implode(", ", array_keys($allConfig)) . "\n";

try {
    $storage = new JsonStorageConnection('localhost', 5555);
    echo "Attempting to connect...\n";
    $storage->connect();
    echo "Connected successfully.\n";
    
    // Test basic PING
    $response = $storage->get("PING");
    echo "PING test: " . ($response === "PONG" ? "PASSED" : "FAILED") . "\n";
    
} catch (Exception $e) {
    echo "\nError: " . $e->getMessage() . "\n";
} finally {
    if (isset($storage)) {
        $storage->close();
    }
}