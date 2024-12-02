<?php

require_once __DIR__ . '/PerformanceTestBase.php';

class GetPerformanceTest extends PerformanceTestBase {
    private $dataTypes = ['string' => 0, 'json' => 0, 'number' => 0, 'complex' => 0];
    
    public function run() {
        try {
            $this->initialize();
            echo "Starting GET Performance Test...\n";
            echo "Retrieving {$this->recordCount} records...\n";
            
            $successCount = 0;
            $batchStart = microtime(true);
            
            for ($i = 0; $i < $this->recordCount; $i++) {
                $key = "perf:test:{$i}";
                $result = $this->storage->get($key);
                
                if ($result !== null) {
                    $successCount++;
                    $this->analyzeDataType($result);
                }
                
                if (($i + 1) % $this->batchSize === 0) {
                    $this->reportBatchProgress($i, $batchStart);
                    $batchStart = microtime(true);
                }
            }
            
            $totalTime = microtime(true) - $this->startTime;
            $peakMemory = memory_get_peak_usage();
            $this->printSummary($totalTime, $successCount, $peakMemory);
            
        } catch (Exception $e) {
            echo "\nERROR: " . $e->getMessage() . "\n";
            echo "Stack trace:\n" . $e->getTraceAsString() . "\n";
        } finally {
            if (isset($this->storage)) {
                $this->storage->close();
            }
        }
    }
    
    private function analyzeDataType($result) {
        if (is_array($result)) {
            if (isset($result['coordinates'])) {
                $this->dataTypes['complex']++;
            } else {
                $this->dataTypes['json']++;
            }
        } elseif (is_numeric($result)) {
            $this->dataTypes['number']++;
        } else {
            $this->dataTypes['string']++;
        }
    }
    
    protected function printSummary($totalTime, $successCount, $peakMemory) {
        parent::printSummary($totalTime, $successCount, $peakMemory);
        
        echo "\nData Type Distribution:\n";
        foreach ($this->dataTypes as $type => $count) {
            $percentage = ($count / $successCount) * 100;
            echo "{$type}: {$count} (" . number_format($percentage, 1) . "%)\n";
        }
    }
}

// Run the test
$test = new GetPerformanceTest();
$test->run();