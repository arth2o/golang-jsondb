const { JsonStorageConnection } = require('./JsonStorageConnection');

console.log('Starting basic commands example...');

async function runBasicCommands() {
    let storage;
    
    try {
        // Initialize connection
        storage = new JsonStorageConnection('localhost', 5555);
        console.log('Connecting to server...');
        await storage.connect();
        
        // Basic commands demonstration
        console.log('\nRunning basic commands:');
        
        // PING test
        console.log('\n=== PING Test ===');
        const pingResponse = await storage.get('PING');
        console.log('PING response:', pingResponse);
        
        // SET command
        console.log('\n=== SET Command ===');
        const key = 'example:test:1';
        const value = 'Hello from basic commands!';
        const success = await storage.set(key, value);
        console.log('SET result:', success ? 'SUCCESS' : 'FAILED');
        
        // GET command
        console.log('\n=== GET Command ===');
        const retrieved = await storage.get(key);
        console.log('GET result:', retrieved);
        
        // SET with TTL
        console.log('\n=== SET with TTL ===');
        const ttlKey = 'example:ttl:1';
        const ttlValue = 'This will expire in 5 seconds';
        await storage.set(ttlKey, ttlValue, 5);
        console.log('Initial TTL value:', await storage.get(ttlKey));
        console.log('TTL remaining:', await storage.ttl(ttlKey), 'seconds');
        
        // Wait and check TTL
        await new Promise(resolve => setTimeout(resolve, 2000));
        console.log('After 2 seconds - TTL remaining:', await storage.ttl(ttlKey), 'seconds');
        
        // DELETE command
        console.log('\n=== DELETE Command ===');
        await storage.del(key);
        const checkDeleted = await storage.get(key);
        console.log('After DELETE - value exists?:', checkDeleted === null ? 'No' : 'Yes');
        
    } catch (error) {
        console.error('\nError:', error.message);
        console.error('Stack trace:\n', error.stack);
    } finally {
        if (storage) {
            storage.close();
            console.log('\nConnection closed.');
        }
    }
}

runBasicCommands();