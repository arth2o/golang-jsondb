<?php

require_once __DIR__ . '/PerformanceTestBase.php';

class JsonPerformanceTest extends PerformanceTestBase {
    private $dataTypes = ['user' => 0, 'product' => 0, 'order' => 0, 'log' => 0, 'metrics' => 0];
    private $setTime;
    private $getTime;
    
    public function run() {
        try {
            $this->initialize();
            echo "Starting JSON Performance Test...\n";
            echo "Generating and inserting {$this->recordCount} random JSON records...\n";
            
            $this->runSetPhase();
            $this->runGetPhase();
            
            $totalTime = microtime(true) - $this->startTime;
            $this->printFinalSummary($totalTime);
            
        } catch (Exception $e) {
            echo "\nERROR: " . $e->getMessage() . "\n";
            echo "Stack trace:\n" . $e->getTraceAsString() . "\n";
        } finally {
            if (isset($this->storage)) {
                $this->storage->close();
            }
        }
    }
    
    private function runSetPhase() {
        $setStartTime = microtime(true);
        $successCount = 0;
        $batchStart = microtime(true);
        
        echo "\nPHASE 1: SET Operations\n";
        for ($i = 0; $i < $this->recordCount; $i++) {
            $data = $this->generateRandomJson();
            $key = "json:perf:test:{$i}";
            
            if ($this->storage->set($key, $data)) {
                $successCount++;
                $this->dataTypes[$data['type']]++;
            }
            
            if (($i + 1) % $this->batchSize === 0) {
                $this->reportBatchProgress($i, $batchStart);
                $batchStart = microtime(true);
            }
        }
        
        $this->setTime = microtime(true) - $setStartTime;
        $this->printPhaseResults("SET", $this->setTime, $successCount);
    }
    
    private function runGetPhase() {
        $getStartTime = microtime(true);
        $successCount = 0;
        $batchStart = microtime(true);
        
        echo "\nPHASE 2: GET Operations\n";
        for ($i = 0; $i < $this->recordCount; $i++) {
            $key = "json:perf:test:{$i}";
            if ($this->storage->get($key) !== null) {
                $successCount++;
            }
            
            if (($i + 1) % $this->batchSize === 0) {
                $this->reportBatchProgress($i, $batchStart);
                $batchStart = microtime(true);
            }
        }
        
        $this->getTime = microtime(true) - $getStartTime;
        $this->printPhaseResults("GET", $this->getTime, $successCount);
    }
    
    private function generateRandomJson() {
        $types = ['user', 'product', 'order', 'log', 'metrics'];
        $type = $types[array_rand($types)];
        
        switch ($type) {
            case 'user':
                return [
                    'type' => 'user',
                    'id' => rand(1000, 9999),
                    'name' => 'User_' . rand(1, 1000),
                    'email' => 'user' . rand(1, 1000) . '@example.com',
                    'preferences' => [
                        'theme' => ['light', 'dark'][array_rand(['light', 'dark'])],
                        'notifications' => (bool)rand(0, 1),
                        'language' => ['en', 'es', 'fr', 'de'][array_rand(['en', 'es', 'fr', 'de'])]
                    ],
                    'metadata' => [
                        'lastLogin' => date('Y-m-d H:i:s', time() - rand(0, 86400 * 30)),
                        'loginCount' => rand(1, 100)
                    ]
                ];
            
            case 'product':
                return [
                    'type' => 'product',
                    'id' => rand(1000, 9999),
                    'sku' => 'PRD-' . rand(100000, 999999),
                    'price' => round(rand(100, 10000) / 100, 2),
                    'inventory' => [
                        'quantity' => rand(0, 1000),
                        'locations' => array_map(function() {
                            return ['id' => 'WH-' . rand(1, 5), 'qty' => rand(0, 200)];
                        }, range(1, rand(1, 3)))
                    ],
                    'attributes' => [
                        'color' => ['red', 'blue', 'green', 'black'][array_rand(['red', 'blue', 'green', 'black'])],
                        'size' => ['S', 'M', 'L', 'XL'][array_rand(['S', 'M', 'L', 'XL'])],
                        'weight' => rand(100, 5000)
                    ]
                ];
            
            case 'order':
                return [
                    'type' => 'order',
                    'id' => 'ORD-' . rand(100000, 999999),
                    'status' => ['pending', 'processing', 'shipped', 'delivered'][array_rand(['pending', 'processing', 'shipped', 'delivered'])],
                    'items' => array_map(function() {
                        return [
                            'productId' => rand(1000, 9999),
                            'quantity' => rand(1, 5),
                            'price' => round(rand(100, 10000) / 100, 2)
                        ];
                    }, range(1, rand(1, 5))),
                    'customer' => [
                        'id' => rand(1000, 9999),
                        'name' => 'Customer_' . rand(1, 1000),
                        'address' => [
                            'street' => rand(100, 999) . ' Main St',
                            'city' => 'City_' . rand(1, 100),
                            'zipCode' => rand(10000, 99999)
                        ]
                    ]
                ];
            
            default:
                return [
                    'type' => $type,
                    'timestamp' => time(),
                    'data' => [
                        'id' => rand(1000, 9999),
                        'value' => rand(1, 1000),
                        'tags' => array_map(function() {
                            return 'tag_' . rand(1, 10);
                        }, range(1, rand(1, 5)))
                    ]
                ];
        }
    }
    
    private function printPhaseResults($phase, $time, $successCount) {
        echo "\n{$phase} Operations:\n";
        echo "  Total Time: " . number_format($time, 2) . " seconds\n";
        echo "  Average Time per {$phase}: " . number_format(($time / $this->recordCount) * 1000, 4) . " ms\n";
        echo "  Success Rate: " . number_format(($successCount / $this->recordCount) * 100, 2) . "%\n";
    }
    
    private function printFinalSummary($totalTime) {
        $peakMemory = memory_get_peak_usage();
        
        echo "\nOverall Statistics:\n";
        echo "Total Time: " . number_format($totalTime, 2) . " seconds\n";
        echo "Initial Memory: " . $this->formatBytes($this->initialMemory) . "\n";
        echo "Peak Memory Usage: " . $this->formatBytes($peakMemory) . "\n";
        echo "Operations per Second: " . number_format(($this->recordCount * 2) / $totalTime, 2) . "\n";
        
        echo "\nData Type Distribution:\n";
        foreach ($this->dataTypes as $type => $count) {
            $percentage = ($count / array_sum($this->dataTypes)) * 100;
            echo "{$type}: {$count} (" . number_format($percentage, 1) . "%)\n";
        }
    }
}

// Run the test
$test = new JsonPerformanceTest();
$test->run();