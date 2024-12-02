<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

echo "Starting JsonStorage connection test...\n";

// Debug: Show all config values
echo "Loading configuration...\n";
$allConfig = ConfigManager::all();
echo "Available configuration keys: " . implode(", ", array_keys($allConfig)) . "\n";

try {
    // Create and connect
    $storage = new JsonStorageConnection('localhost', 5555);
    $storage->connect();
    
    // Run all test suites
    runBasicTests($storage);
    runDataTypeTests($storage);
    runTTLTests($storage);
    runComplexDataTests($storage);
    
} catch (Exception $e) {
    echo "\nError: " . $e->getMessage() . "\n";
} 

function runBasicTests($storage) {
    echo "\nRunning Basic Functionality Tests:\n";
    
    // Test PING
    $response = $storage->get("PING");
    echo "PING test: " . ($response === "PONG" ? "PASSED" : "FAILED") . "\n";

    // Test basic SET/GET
    $storage->set("test:basic:1", "Hello World");
    $response = $storage->get("test:basic:1");
    echo "Basic SET/GET test: " . ($response === "Hello World" ? "PASSED" : "FAILED") . "\n";

    // Test DELETE
    $deleted = $storage->del("test:basic:1");
    $response = $storage->get("test:basic:1");
    echo "DELETE test: " . ($response === null ? "PASSED" : "FAILED") . "\n";

    // Test non-existent key
    $response = $storage->get("test:nonexistent");
    echo "Non-existent key test: " . ($response === null ? "PASSED" : "FAILED") . "\n";
}

function runDataTypeTests($storage) {
    echo "\nRunning Data Type Tests:\n";
    
    $testCases = [
        'test:string' => [
            'value' => "Simple string test",
            'type' => 'string'
        ],
        'test:integer' => [
            'value' => 42,
            'type' => 'integer'
        ],
        'test:float' => [
            'value' => 3.14159,
            'type' => 'double'
        ],
        'test:boolean' => [
            'value' => true,
            'type' => 'boolean'
        ],
        'test:null' => [
            'value' => null,
            'type' => 'NULL'
        ],
        'test:special_chars' => [
            'value' => "Special chars: !@#$%^&*()",
            'type' => 'string'
        ]
    ];
    
    foreach ($testCases as $key => $test) {
        $storage->set($key, $test['value']);
        $retrieved = $storage->get($key);
        
        $typeMatch = gettype($retrieved) === $test['type'];
        $valueMatch = $retrieved === $test['value'];
        
        echo sprintf(
            "%s: %s (Type: %s, Expected: %s) - %s\n",
            $key,
            ($typeMatch && $valueMatch) ? "PASSED" : "FAILED",
            gettype($retrieved),
            $test['type'],
            $test['value'] ?? 'null'
        );
        
        $storage->del($key);
    }
}

function runTTLTests($storage) {
    echo "\nRunning TTL and Expiration Tests:\n";
    
    // Test TTL setting
    $storage->set("test:ttl:1", "Expires in 2 seconds", 2);
    $ttl = $storage->ttl("test:ttl:1");
    echo "TTL test (should be ~2): " . ($ttl <= 2 && $ttl > 0 ? "PASSED" : "FAILED") . " (TTL: $ttl)\n";
    
    // Test expiration
    echo "Waiting for key to expire...\n";
    sleep(3);
    $expired = $storage->get("test:ttl:1");
    echo "Expiration test: " . ($expired === null ? "PASSED" : "FAILED") . "\n";
    
    // Test no expiration
    $storage->set("test:ttl:2", "No expiration", -1);
    $ttl = $storage->ttl("test:ttl:2");
    echo "No expiration test: " . ($ttl === -1 ? "PASSED" : "FAILED") . "\n";
    $storage->del("test:ttl:2");
}

function runComplexDataTests($storage) {
    echo "\nRunning Complex Data Structure Tests:\n";
    
    $testCases = [
        'test:array:simple' => [
            'data' => ['apple', 'banana', 'orange'],
            'desc' => 'Simple array'
        ],
        'test:array:assoc' => [
            'data' => ['name' => 'John', 'age' => 30, 'city' => 'New York'],
            'desc' => 'Associative array'
        ],
        'test:nested:deep' => [
            'data' => [
                'user' => [
                    'profile' => [
                        'name' => 'Jane Doe',
                        'settings' => [
                            'theme' => 'dark',
                            'notifications' => true,
                            'preferences' => [
                                'language' => 'en',
                                'timezone' => 'UTC'
                            ]
                        ]
                    ]
                ]
            ],
            'desc' => 'Deeply nested structure'
        ],
        'test:mixed:types' => [
            'data' => [
                'string' => 'text',
                'number' => 42,
                'float' => 3.14,
                'boolean' => true,
                'null' => null,
                'array' => [1, 2, 3],
                'object' => ['key' => 'value']
            ],
            'desc' => 'Mixed data types'
        ]
    ];
    
    foreach ($testCases as $key => $test) {
        echo "\nTesting {$test['desc']}:\n";
        
        $storage->set($key, $test['data']);
        $retrieved = $storage->get($key);
        
        $matches = json_encode($test['data']) === json_encode($retrieved);
        echo "Data integrity: " . ($matches ? "PASSED" : "FAILED") . "\n";
        
        if (!$matches) {
            echo "Original: " . json_encode($test['data']) . "\n";
            echo "Retrieved: " . json_encode($retrieved) . "\n";
        }
        
        $storage->del($key);
    }
}