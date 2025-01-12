// Products module
const products = {
    async getAll() {
        return await api.get('/products');
    },

    async create(product) {
        return await api.post('/products', product);
    }
};
