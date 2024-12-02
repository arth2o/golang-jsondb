<?php

require_once __DIR__ . '/JsonStorageConnectionClass.php';

abstract class PerformanceTestBase {
    protected $storage;
    protected $startTime;
    protected $initialMemory;
    protected $recordCount = 100000;
    protected $batchSize = 1000;
    
    protected function initialize() {
        echo "Loading configuration...\n";
        $allConfig = ConfigManager::all();
        echo "Available configuration keys: " . implode(", ", array_keys($allConfig)) . "\n";
        
        $this->startTime = microtime(true);
        $this->initialMemory = memory_get_usage();
        
        $this->storage = new JsonStorageConnection('localhost', 5555);
        $this->storage->connect();
    }
    
    protected function formatBytes($bytes, $precision = 2) {
        $units = ['B', 'KB', 'MB', 'GB'];
        $bytes = max($bytes, 0);
        $pow = floor(($bytes ? log($bytes) : 0) / log(1024));
        $pow = min($pow, count($units) - 1);
        return round($bytes / pow(1024, $pow), $precision) . ' ' . $units[$pow];
    }
    
    protected function printSummary($totalTime, $successCount, $peakMemory) {
        echo "\nPerformance Summary:\n";
        echo "Total Time: " . number_format($totalTime, 2) . " seconds\n";
        echo "Average Time per Record: " . number_format(($totalTime / $this->recordCount) * 1000, 4) . " ms\n";
        echo "Initial Memory: " . $this->formatBytes($this->initialMemory) . "\n";
        echo "Peak Memory Usage: " . $this->formatBytes($peakMemory) . "\n";
        echo "Successful Records: {$successCount}/{$this->recordCount}\n";
        echo "Operations per Second: " . number_format($this->recordCount / $totalTime, 2) . "\n";
    }
    
    protected function reportBatchProgress($i, $batchStart) {
        $batchTime = microtime(true) - $batchStart;
        $currentMemory = memory_get_usage();
        echo "Batch " . (($i + 1) / $this->batchSize) . " complete. ";
        echo "Time: " . number_format($batchTime, 2) . "s, ";
        echo "Memory: " . $this->formatBytes($currentMemory) . "\n";
    }
    
    abstract public function run();
}