<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

echo "Starting basic commands example...\n";

try {
    // Initialize connection
    $storage = new JsonStorageConnection('localhost', 5555);
    echo "Connecting to server...\n";
    $storage->connect();
    
    // Basic commands demonstration
    echo "\nRunning basic commands:\n";
    
    // PING test
    echo "\n=== PING Test ===\n";
    $response = $storage->get("PING");
    echo "PING response: " . $response . "\n";
    
    // SET command
    echo "\n=== SET Command ===\n";
    $key = "example:test:1";
    $value = "Hello from basic commands!";
    $success = $storage->set($key, $value);
    echo "SET result: " . ($success ? "SUCCESS" : "FAILED") . "\n";
    
    // GET command
    echo "\n=== GET Command ===\n";
    $retrieved = $storage->get($key);
    echo "GET result: " . $retrieved . "\n";
    
    // SET with TTL
    echo "\n=== SET with TTL ===\n";
    $ttlKey = "example:ttl:1";
    $ttlValue = "This will expire in 5 seconds";
    $storage->set($ttlKey, $ttlValue, 5);
    echo "Initial TTL value: " . $storage->get($ttlKey) . "\n";
    echo "TTL remaining: " . $storage->ttl($ttlKey) . " seconds\n";
    
    // Wait and check TTL
    sleep(2);
    echo "After 2 seconds - TTL remaining: " . $storage->ttl($ttlKey) . " seconds\n";
    
    // DELETE command
    echo "\n=== DELETE Command ===\n";
    $storage->del($key);
    $checkDeleted = $storage->get($key);
    echo "After DELETE - value exists?: " . ($checkDeleted === null ? "No" : "Yes") . "\n";
    
} catch (Exception $e) {
    echo "\nError: " . $e->getMessage() . "\n";
    echo "Stack trace:\n" . $e->getTraceAsString() . "\n";
} finally {
    if (isset($storage)) {
        $storage->close();
        echo "\nConnection closed.\n";
    }
}