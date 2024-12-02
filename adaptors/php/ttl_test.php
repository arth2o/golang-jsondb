<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

echo "Starting TTL Test...\n";

// Debug: Show all config values
echo "Loading configuration...\n";
$allConfig = ConfigManager::all();
echo "Available configuration keys: " . implode(", ", array_keys($allConfig)) . "\n";

try {
    $storage = new JsonStorageConnection('localhost', 5555);
    $storage->connect();
    
    // Run TTL test suites
    runBasicTTLTests($storage);
    runComplexTTLTests($storage);
    runExpirationTests($storage);
    
} catch (Exception $e) {
    echo "\nError: " . $e->getMessage() . "\n";
} finally {
    if (isset($storage)) {
        $storage->close();
    }
}

function runBasicTTLTests($storage) {
    echo "\nRunning Basic TTL Tests:\n";
    
    $testCases = [
        [
            'name' => 'No TTL',
            'key' => 'ttl:test:1',
            'value' => 'Permanent Value',
            'ttl' => -1
        ],
        [
            'name' => '5 Second TTL',
            'key' => 'ttl:test:2',
            'value' => 'Temporary Value',
            'ttl' => 5
        ]
    ];
    
    foreach ($testCases as $test) {
        echo "\n=== Testing {$test['name']} ===\n";
        runSingleTTLTest($storage, $test);
    }
}

function runComplexTTLTests($storage) {
    echo "\nRunning Complex Data TTL Tests:\n";
    
    $testCases = [
        [
            'name' => 'JSON with TTL',
            'key' => 'ttl:test:json',
            'value' => ['name' => 'John', 'expires' => true],
            'ttl' => 10
        ],
        [
            'name' => 'Array with TTL',
            'key' => 'ttl:test:array',
            'value' => [1, 2, 3, 4, 5],
            'ttl' => 8
        ]
    ];
    
    foreach ($testCases as $test) {
        echo "\n=== Testing {$test['name']} ===\n";
        runSingleTTLTest($storage, $test);
    }
}

function runExpirationTests($storage) {
    echo "\n=== Testing Expiration ===\n";
    $expireKey = 'ttl:expire:test';
    $storage->del($expireKey);
    
    echo "Setting value with 3s TTL...\n";
    $storage->set($expireKey, 'This will expire', 3);
    
    $initialValue = $storage->get($expireKey);
    if ($initialValue === null) {
        echo "ERROR: Failed to set initial value\n";
        return;
    }
    
    echo "Initial TTL: " . $storage->ttl($expireKey) . " seconds\n";
    echo "Initial value: " . ($initialValue ?? 'null') . "\n";
    
    echo "Waiting for expiration (4 seconds)...\n";
    sleep(4);
    
    $finalTTL = $storage->ttl($expireKey);
    echo "TTL after wait: " . ($finalTTL === -2 ? "Key expired" : $finalTTL . " seconds") . "\n";
    
    $finalValue = $storage->get($expireKey);
    echo "Final value: " . ($finalValue === null ? "Expired (null)" : $finalValue) . "\n";
}

function runSingleTTLTest($storage, $test) {
    $storage->del($test['key']);
    
    echo "Setting value with TTL {$test['ttl']}...\n";
    $success = $storage->set($test['key'], $test['value'], $test['ttl']);
    
    $verifyValue = $storage->get($test['key']);
    $success = ($verifyValue !== null);
    
    echo "Set operation: " . ($success ? "PASSED" : "FAILED") . "\n";
    
    $initialTTL = $storage->ttl($test['key']);
    echo "Initial TTL: " . ($initialTTL === -1 ? "No expiration" : $initialTTL . " seconds") . "\n";
    
    if ($test['ttl'] > 0) {
        sleep(2);
        $remainingTTL = $storage->ttl($test['key']);
        echo "TTL after 2s: " . ($remainingTTL === -1 ? "No expiration" : $remainingTTL . " seconds") . "\n";
    }
    
    $storage->del($test['key']);
}