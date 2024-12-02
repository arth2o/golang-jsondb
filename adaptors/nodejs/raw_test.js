const { JsonStorageConnection } = require('./JsonStorageConnection');
const { ConfigManager } = require('./config/ConfigManager');

console.log('Starting JsonStorage connection test...');

// Debug: Show all config values
console.log('Loading configuration...');
const allConfig = ConfigManager.all();
console.log('Available configuration keys:', Object.keys(allConfig).join(', '));

async function runBasicTests(storage) {
    console.log('\nRunning Basic Functionality Tests:');
    
    // Test PING
    const pingResponse = await storage.get('PING');
    console.log(`PING test: ${pingResponse === 'PONG' ? 'PASSED' : 'FAILED'}`);

    // Test basic SET/GET
    await storage.set('test:basic:1', 'Hello World');
    const getResponse = await storage.get('test:basic:1');
    console.log(`Basic SET/GET test: ${getResponse === 'Hello World' ? 'PASSED' : 'FAILED'}`);

    // Test DELETE
    await storage.del('test:basic:1');
    const deletedResponse = await storage.get('test:basic:1');
    console.log(`DELETE test: ${deletedResponse === null ? 'PASSED' : 'FAILED'}`);

    // Test non-existent key
    const nonExistentResponse = await storage.get('test:nonexistent');
    console.log(`Non-existent key test: ${nonExistentResponse === null ? 'PASSED' : 'FAILED'}`);
}

async function runDataTypeTests(storage) {
    console.log('\nRunning Data Type Tests:');
    
    const testCases = {
        'test:string': {
            value: 'Simple string test',
            type: 'string'
        },
        'test:integer': {
            value: 42,
            type: 'number'
        },
        'test:float': {
            value: 3.14159,
            type: 'number'
        },
        'test:boolean': {
            value: true,
            type: 'boolean'
        },
        'test:null': {
            value: null,
            type: 'object'
        },
        'test:special_chars': {
            value: 'Special chars: !@#$%^&*()',
            type: 'string',
            compare: (original, retrieved) => {
                const normalizedOriginal = original.normalize();
                const normalizedRetrieved = retrieved
                    .replace(/\\u([0-9a-fA-F]{4})/g, (_, code) => String.fromCharCode(parseInt(code, 16)))
                    .normalize();
                return normalizedOriginal === normalizedRetrieved;
            }
        }
    };
    
    for (const [key, test] of Object.entries(testCases)) {
        await storage.set(key, test.value);
        const retrieved = await storage.get(key);
        
        const typeMatch = typeof retrieved === test.type;
        const valueMatch = test.compare ? 
            test.compare(test.value, retrieved) : 
            retrieved === test.value;
        
        console.log(
            `${key}: ${typeMatch && valueMatch ? 'PASSED' : 'FAILED'} ` +
            `(Type: ${typeof retrieved}, Expected: ${test.type}) - ${test.value}`
        );
        
        await storage.del(key);
    }
}

async function runTTLTests(storage) {
    console.log('\nRunning TTL and Expiration Tests:');
    
    // Test TTL setting
    await storage.set('test:ttl:1', 'Expires in 2 seconds', 2);
    const ttl = await storage.ttl('test:ttl:1');
    console.log(`TTL test (should be ~2): ${ttl <= 2 && ttl > 0 ? 'PASSED' : 'FAILED'} (TTL: ${ttl})`);
    
    // Test expiration
    console.log('Waiting for key to expire...');
    await new Promise(resolve => setTimeout(resolve, 3000));
    const expired = await storage.get('test:ttl:1');
    console.log(`Expiration test: ${expired === null ? 'PASSED' : 'FAILED'}`);
    
    // Test no expiration
    await storage.set('test:ttl:2', 'No expiration', -1);
    const noExpTtl = await storage.ttl('test:ttl:2');
    console.log(`No expiration test: ${noExpTtl === -1 ? 'PASSED' : 'FAILED'}`);
    await storage.del('test:ttl:2');
}

async function runComplexDataTests(storage) {
    console.log('\nRunning Complex Data Structure Tests:');
    
    const testCases = {
        'test:array:simple': {
            data: ['apple', 'banana', 'orange'],
            desc: 'Simple array'
        },
        'test:array:assoc': {
            data: { name: 'John', age: 30, city: 'New York' },
            desc: 'Associative array'
        },
        'test:nested:deep': {
            data: {
                user: {
                    profile: {
                        name: 'Jane Doe',
                        settings: {
                            theme: 'dark',
                            notifications: true,
                            preferences: {
                                language: 'en',
                                timezone: 'UTC'
                            }
                        }
                    }
                }
            },
            desc: 'Deeply nested structure'
        },
        'test:mixed:types': {
            data: {
                string: 'text',
                number: 42,
                float: 3.14,
                boolean: true,
                null: null,
                array: [1, 2, 3],
                object: { key: 'value' }
            },
            desc: 'Mixed data types'
        }
    };
    
    for (const [key, test] of Object.entries(testCases)) {
        console.log(`\nTesting ${test.desc}:`);
        
        await storage.set(key, test.data);
        const retrieved = await storage.get(key);
        
        const matches = JSON.stringify(test.data) === JSON.stringify(retrieved);
        console.log(`Data integrity: ${matches ? 'PASSED' : 'FAILED'}`);
        
        if (!matches) {
            console.log('Original:', JSON.stringify(test.data));
            console.log('Retrieved:', JSON.stringify(retrieved));
        }
        
        await storage.del(key);
    }
}

async function runTests() {
    try {
        const storage = new JsonStorageConnection('localhost', 5555);
        await storage.connect();
        
        await runBasicTests(storage);
        await runDataTypeTests(storage);
        await runTTLTests(storage);
        await runComplexDataTests(storage);
        
        storage.close();
    } catch (error) {
        console.error('\nError:', error.message);
        console.error(error.stack);
    }
}

runTests();