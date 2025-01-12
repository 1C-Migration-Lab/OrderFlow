// API client for interacting with the backend
const API = {
    baseUrl: 'http://localhost:8080/api',

    // Generic request method
    async request(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.baseUrl}${endpoint}`, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers,
                },
            });

            if (!response.ok) {
                const error = await response.json().catch(() => ({}));
                throw new Error(error.error || `HTTP error! status: ${response.status}`);
            }

            // Для DELETE запросов может не быть тела ответа
            if (options.method === 'DELETE') {
                return null;
            }

            return await response.json().catch(() => null);
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    },

    // Clients
    async getClients() {
        return this.request('/clients');
    },

    async createClient(client) {
        return this.request('/clients', {
            method: 'POST',
            body: JSON.stringify(client),
        });
    },

    async updateClient(id, client) {
        return this.request(`/clients/${id}`, {
            method: 'PUT',
            body: JSON.stringify(client),
        });
    },

    async deleteClient(id) {
        return this.request(`/clients/${id}`, {
            method: 'DELETE',
        });
    },

    // Products
    async getProducts() {
        return this.request('/products');
    },

    async createProduct(product) {
        return this.request('/products', {
            method: 'POST',
            body: JSON.stringify(product),
        });
    },

    async updateProduct(id, product) {
        return this.request(`/products/${id}`, {
            method: 'PUT',
            body: JSON.stringify(product),
        });
    },

    async deleteProduct(id) {
        return this.request(`/products/${id}`, {
            method: 'DELETE',
        });
    },

    // Orders
    async getOrders() {
        return this.request('/orders');
    },

    async createOrder(order) {
        return this.request('/orders', {
            method: 'POST',
            body: JSON.stringify(order),
        });
    },

    async updateOrder(id, order) {
        return this.request(`/orders/${id}`, {
            method: 'PUT',
            body: JSON.stringify(order),
        });
    },

    async deleteOrder(id) {
        return this.request(`/orders/${id}`, {
            method: 'DELETE',
        });
    },

    async confirmOrder(id) {
        return this.request(`/orders/${id}/confirm`, {
            method: 'POST',
        });
    },

    // Orders by client
    async getOrdersByClient() {
        return this.request('/orders-by-client');
    },

    renderAll() {
        this.showLoading();
        try {
            this.renderClients();
            this.renderProducts();
            this.renderOrders();
        } finally {
            this.hideLoading();
        }
    }
};
