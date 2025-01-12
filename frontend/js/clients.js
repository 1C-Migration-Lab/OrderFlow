// Clients module
const clients = {
    async getAll() {
        return await api.get('/clients');
    },

    async create(client) {
        return await api.post('/clients', client);
    }
};
