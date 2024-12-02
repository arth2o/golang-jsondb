#!/bin/bash

# Set colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

echo -e "${BLUE}Starting PHP Tests...${NC}\n"

# Array of test files
declare -a test_files=(
    "raw_test.php"
    "ttl_test.php"    
    "index.php"    
    "test_connection.php"
    "performance_json_test.php"
    "performance_get_test.php"
    "performance_set_test.php"
    "interactive_test.php"
)

# Function to run test and check result
run_test() {
    local test_file=$1
    echo -e "${BLUE}Running test: ${test_file}${NC}"
    echo "----------------------------------------"
    
    if php -f "$test_file"; then
        echo -e "${GREEN}✓ Test completed successfully${NC}"
    else
        echo -e "${RED}✗ Test failed with exit code $?${NC}"
        exit 1
    fi
    echo "----------------------------------------\n"
}

# Run each test
for test_file in "${test_files[@]}"; do
    if [ -f "$test_file" ]; then
        run_test "$test_file"
    else
        echo -e "${RED}Error: Test file not found: ${test_file}${NC}"
        exit 1
    fi
done

echo -e "${GREEN}All tests completed!${NC}"