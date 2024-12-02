<?php

require_once 'JsonStorageConnectionClass.php';

error_reporting(E_ALL);
ini_set('display_errors', 1);

try {
    echo "Creating JsonStorageConnection instance...\n";
    $storage = new JsonStorageConnection('localhost', 5555);
    
    echo "Attempting to connect...\n";
    $storage->connect();
    echo "Connected successfully!\n";

    // Test cases
    $tests = [
        [
            'name' => 'Simple String',
            'key' => 'test:string',
            'value' => 'Hello World'
        ],
        [
            'name' => 'JSON Object',
            'key' => 'test:json',
            'value' => ['name' => 'John', 'age' => 30]
        ],
        [
            'name' => 'Numbers',
            'key' => 'test:number',
            'value' => 42
        ]
    ];

    foreach ($tests as $test) {
        echo "\n=== Testing {$test['name']} ===\n";
        
        echo "Setting value...\n";
        $storage->set($test['key'], $test['value']);
        
        echo "Getting value...\n";
        $result = $storage->get($test['key']);
        
        echo "Original:\n";
        var_dump($test['value']);
        echo "Retrieved:\n";
        var_dump($result);
        
        $matches = $test['value'] === $result;
        echo "Match: " . ($matches ? "Yes" : "No") . "\n";
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