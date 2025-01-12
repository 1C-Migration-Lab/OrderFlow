// Main application logic
document.addEventListener('DOMContentLoaded', () => {
    // Application state
    const state = {
        clients: [],
        products: [],
        orders: [],
        ordersByClient: [],
        currentSection: 'clients'
    };

    // Initialize application
    const app = {
        async init() {
            await this.loadInitialData();
            this.setupEventListeners();
            this.showSection(state.currentSection);
        },

        async loadInitialData() {
            try {
                const [clients, products, orders, ordersByClient] = await Promise.all([
                    API.getClients(),
                    API.getProducts(),
                    API.getOrders(),
                    API.getOrdersByClient()
                ]);

                state.clients = clients;
                state.products = products;
                state.orders = orders;
                state.ordersByClient = ordersByClient;

                this.renderAll();
            } catch (error) {
                console.error('Failed to load data:', error);
                alert('Failed to load data. Please try again later.');
            }
        },

        setupEventListeners() {
            // Form submissions
            document.getElementById('client-form').addEventListener('submit', this.handleClientSubmit.bind(this));
            document.getElementById('product-form').addEventListener('submit', this.handleProductSubmit.bind(this));
            document.getElementById('order-form').addEventListener('submit', this.handleOrderSubmit.bind(this));
        },

        // UI Helpers
        showSection(sectionId) {
            state.currentSection = sectionId;
            document.querySelectorAll('.section').forEach(section => {
                section.classList.toggle('active', section.id === sectionId);
            });
        },

        showModal(modalId) {
            document.getElementById(modalId).style.display = 'block';
        },

        hideModal(modalId) {
            document.getElementById(modalId).style.display = 'none';
            // Reset form
            document.getElementById(modalId).querySelector('form').reset();
        },

        // Rendering
        renderAll() {
            this.renderClients();
            this.renderProducts();
            this.renderOrders();
        },

        renderClients() {
            const list = document.getElementById('clients-list');
            list.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>INN</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${state.clients.map(client => `
                            <tr>
                                <td>${client.name}</td>
                                <td>${client.inn || '-'}</td>
                                <td>
                                    <button onclick="app.editClient(${client.id})">Edit</button>
                                    <button class="delete" onclick="app.deleteClient(${client.id})">Delete</button>
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        },

        renderProducts() {
            const list = document.getElementById('products-list');
            list.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Unit</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${state.products.map(product => `
                            <tr>
                                <td>${product.name}</td>
                                <td>${product.unit}</td>
                                <td>
                                    <button onclick="app.editProduct(${product.id})">Edit</button>
                                    <button class="delete" onclick="app.deleteProduct(${product.id})">Delete</button>
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        },

        renderOrders() {
            const list = document.getElementById('orders-list');
            list.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>Number</th>
                            <th>Client</th>
                            <th>Date</th>
                            <th>Total Amount</th>
                            <th>Status</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${state.orders.map(order => `
                            <tr>
                                <td>${order.number}</td>
                                <td>${state.clients.find(c => c.id === order.client_id)?.name || 'Unknown'}</td>
                                <td>${new Date(order.date).toLocaleDateString()}</td>
                                <td>${order.total_amount}</td>
                                <td>${order.is_confirmed ? 'Confirmed' : 'Draft'}</td>
                                <td>
                                    ${!order.is_confirmed ? `
                                        <button onclick="app.editOrder(${order.id})">Edit</button>
                                        <button onclick="app.confirmOrder(${order.id})">Confirm</button>
                                    ` : ''}
                                    <button class="delete" onclick="app.deleteOrder(${order.id})">Delete</button>
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        },

        // Event Handlers
        async handleClientSubmit(event) {
            event.preventDefault();
            const form = event.target;
            const formData = new FormData(form);
            
            try {
                const client = {
                    name: formData.get('name'),
                    inn: formData.get('inn')
                };

                const result = await API.createClient(client);
                state.clients.push(result);
                this.renderClients();
                this.hideModal('client-modal');
            } catch (error) {
                console.error('Failed to create client:', error);
                alert('Failed to create client. Please try again.');
            }
        },

        async handleProductSubmit(event) {
            event.preventDefault();
            const form = event.target;
            const formData = new FormData(form);
            
            try {
                const product = {
                    name: formData.get('name'),
                    unit: formData.get('unit')
                };

                const result = await API.createProduct(product);
                state.products.push(result);
                this.renderProducts();
                this.hideModal('product-modal');
            } catch (error) {
                console.error('Failed to create product:', error);
                alert('Failed to create product. Please try again.');
            }
        },

        async handleOrderSubmit(event) {
            event.preventDefault();
            const form = event.target;
            const formData = new FormData(form);
            
            try {
                const order = {
                    client_id: parseInt(formData.get('client_id')),
                    number: formData.get('number'),
                    items: Array.from(document.querySelectorAll('.order-item')).map(item => ({
                        product_id: parseInt(item.querySelector('[name="product_id"]').value),
                        quantity: parseFloat(item.querySelector('[name="quantity"]').value),
                        price: parseFloat(item.querySelector('[name="price"]').value)
                    }))
                };

                const result = await API.createOrder(order);
                state.orders.push(result);
                this.renderOrders();
                this.hideModal('order-modal');
            } catch (error) {
                console.error('Failed to create order:', error);
                alert('Failed to create order. Please try again.');
            }
        },

        // CRUD operations
        async deleteClient(id) {
            if (!confirm('Are you sure you want to delete this client?')) return;
            
            try {
                await API.deleteClient(id);
                state.clients = state.clients.filter(c => c.id !== id);
                this.renderClients();
            } catch (error) {
                console.error('Failed to delete client:', error);
                alert('Failed to delete client. Please try again.');
            }
        },

        async deleteProduct(id) {
            if (!confirm('Are you sure you want to delete this product?')) return;
            
            try {
                await API.deleteProduct(id);
                state.products = state.products.filter(p => p.id !== id);
                this.renderProducts();
            } catch (error) {
                console.error('Failed to delete product:', error);
                alert('Failed to delete product. Please try again.');
            }
        },

        async deleteOrder(id) {
            if (!confirm('Are you sure you want to delete this order?')) return;
            
            try {
                await API.deleteOrder(id);
                state.orders = state.orders.filter(o => o.id !== id);
                this.renderOrders();
            } catch (error) {
                console.error('Failed to delete order:', error);
                alert('Failed to delete order. Please try again.');
            }
        },

        async confirmOrder(id) {
            try {
                const result = await API.confirmOrder(id);
                const index = state.orders.findIndex(o => o.id === id);
                if (index !== -1) {
                    state.orders[index] = result;
                }
                this.renderOrders();
            } catch (error) {
                console.error('Failed to confirm order:', error);
                alert('Failed to confirm order. Please try again.');
            }
        },

        // Order items management
        addOrderItem() {
            const container = document.getElementById('order-items');
            const itemDiv = document.createElement('div');
            itemDiv.className = 'order-item';
            itemDiv.innerHTML = `
                <select name="product_id" required onchange="app.updateOrderTotal()">
                    <option value="">Select Product</option>
                    ${state.products.map(p => `
                        <option value="${p.id}">${p.name} (${p.unit})</option>
                    `).join('')}
                </select>
                <input type="number" name="quantity" placeholder="Quantity" required min="0" step="0.001" onchange="app.updateOrderTotal()">
                <input type="number" name="price" placeholder="Price" required min="0" step="0.01" onchange="app.updateOrderTotal()">
                <button type="button" class="delete" onclick="this.parentElement.remove(); app.updateOrderTotal()">Remove</button>
            `;
            container.appendChild(itemDiv);
        },

        updateOrderTotal() {
            const items = document.querySelectorAll('.order-item');
            const total = Array.from(items).reduce((sum, item) => {
                const quantity = parseFloat(item.querySelector('[name="quantity"]').value) || 0;
                const price = parseFloat(item.querySelector('[name="price"]').value) || 0;
                return sum + (quantity * price);
            }, 0);
            document.getElementById('order-total').textContent = total.toFixed(2);
        }
    };

    // Make app globally available
    window.app = app;

    // Start the application
    app.init();
});
