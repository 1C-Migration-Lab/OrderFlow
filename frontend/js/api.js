// API endpoints configuration
const API_BASE_URL = 'http://localhost:8080/api';

// API client implementation
class ApiClient {
    async get(endpoint) {
        const response = await fetch(`${API_BASE_URL}${endpoint}`);
        return response.json();
    }

    async post(endpoint, data) {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return response.json();
    }
}

const api = new ApiClient();
