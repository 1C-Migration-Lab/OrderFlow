// Orders module
const orders = {
    async getAll() {
        return await api.get('/orders');
    },

    async create(order) {
        return await api.post('/orders', order);
    },

    async confirm(orderId) {
        return await api.post(`/orders/${orderId}/confirm`);
    },

    async getByClient() {
        return await api.get('/orders-by-client');
    }
};
