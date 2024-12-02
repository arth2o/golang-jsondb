const fs = require('fs');
const path = require('path');

class EnvLoader {
    static load(environment) {
        const envFile = path.join(__dirname, '../../../jsondb/.env.' + environment);

        if (!fs.existsSync(envFile)) {
            throw new Error(`Environment file not found: ${envFile}`);
        }

        const config = {};
        const content = fs.readFileSync(envFile, 'utf8');
        const lines = content.split('\n');

        for (const line of lines) {
            const trimmedLine = line.trim();
            if (!trimmedLine || trimmedLine.startsWith('#')) continue;

            const [key, ...valueParts] = trimmedLine.split('=');
            let value = valueParts.join('=').trim();

            // Remove surrounding quotes if present
            if ((value.startsWith("'") && value.endsWith("'")) ||
                (value.startsWith('"') && value.endsWith('"'))) {
                value = value.slice(1, -1);
            }

            config[key.trim()] = value;
        }

        return config;
    }
}

module.exports = { EnvLoader };