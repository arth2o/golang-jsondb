<?php

require_once __DIR__ . '/PerformanceTestBase.php';

class SetPerformanceTest extends PerformanceTestBase {
    public function run() {
        try {
            $this->initialize();
            echo "Starting SET Performance Test...\n";
            
            $successCount = 0;
            $batchStart = microtime(true);
            
            // Create test data array
            $testData = [
                ['type' => 'string', 'value' => 'Hello World'],
                ['type' => 'json', 'value' => ['name' => 'John', 'age' => 30]],
                ['type' => 'number', 'value' => 42],
                ['type' => 'complex', 'value' => [
                    'id' => 1,
                    'data' => [
                        'coordinates' => [10.5, 20.3],
                        'active' => true,
                        'tags' => ['test', 'performance', 'benchmark']
                    ]
                ]]
            ];
            
            // Main test loop
            for ($i = 0; $i < $this->recordCount; $i++) {
                $data = $testData[$i % count($testData)];
                $key = "perf:test:{$i}";
                
                if ($this->storage->set($key, $data['value'])) {
                    $successCount++;
                }
                
                // Batch reporting
                if (($i + 1) % $this->batchSize === 0) {
                    $this->reportBatchProgress($i, $batchStart);
                    $batchStart = microtime(true);
                }
            }
            
            $totalTime = microtime(true) - $this->startTime;
            $this->printSummary($totalTime, $successCount, memory_get_peak_usage());
            
        } catch (Exception $e) {
            echo "\nERROR: " . $e->getMessage() . "\n";
            echo "Stack trace:\n" . $e->getTraceAsString() . "\n";
        } finally {
            if (isset($this->storage)) {
                $this->storage->close();
            }
        }
    }
}

// Run the test
$test = new SetPerformanceTest();
$test->run();