/**
 * Common JavaScript utilities for CRM Stomatology
 */

// Toast Notification System
const Toast = {
    show(message, type = 'info', duration = 3000) {
        const toast = document.getElementById('toast');
        if (!toast) return;

        toast.textContent = message;
        toast.className = `show ${type}`;

        setTimeout(() => {
            toast.className = '';
        }, duration);
    },

    success(message) {
        this.show(message, 'success');
    },

    error(message) {
        this.show(message, 'error');
    },

    info(message) {
        this.show(message, 'info');
    }
};

// For backward compatibility
function showNotification(message, type = 'info') {
    Toast.show(message, type);
}

// Auth utilities
const Auth = {
    check() {
        const userRole = localStorage.getItem('userRole');
        if (!userRole) {
            window.location.href = '/login.html';
            return false;
        }
        return true;
    },

    logout() {
        localStorage.clear();
        window.location.href = '/login.html';
    },

    getUser() {
        return {
            name: localStorage.getItem('userDisplayName') || 'User',
            role: localStorage.getItem('userRole') || 'user',
            isAdmin: localStorage.getItem('userRole') === 'admin'
        };
    },

    setupLogoutButton() {
        const nav = document.querySelector('.nav');
        if (!nav) return;

        const logoutBtn = document.createElement('a');
        logoutBtn.href = '#';
        logoutBtn.innerHTML = '🚪 Выйти';
        logoutBtn.style.backgroundColor = '#dc3545';
        logoutBtn.onclick = (e) => {
            e.preventDefault();
            this.logout();
        };
        nav.appendChild(logoutBtn);
    }
};

// Date formatting utilities
const DateUtils = {
    format(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '-';
        return date.toLocaleDateString('ru-RU');
    },

    formatTime(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '-';
        return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    },

    formatDateTime(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '-';
        return date.toLocaleString('ru-RU');
    },

    toInputFormat(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '';
        return date.toISOString().split('T')[0];
    }
};

// Currency formatting
const Currency = {
    format(amount) {
        if (amount === null || amount === undefined) return '0';
        return new Intl.NumberFormat('ru-RU').format(amount);
    },

    formatWithSymbol(amount) {
        return `${this.format(amount)} ₸`;
    }
};

// Phone formatting
const Phone = {
    format(phone) {
        if (!phone) return '';
        // Remove all non-digits
        const digits = phone.replace(/\D/g, '');

        // Format as +7 (XXX) XXX-XX-XX
        if (digits.length === 11 && digits[0] === '7') {
            return `+7 (${digits.slice(1,4)}) ${digits.slice(4,7)}-${digits.slice(7,9)}-${digits.slice(9,11)}`;
        }
        if (digits.length === 10) {
            return `+7 (${digits.slice(0,3)}) ${digits.slice(3,6)}-${digits.slice(6,8)}-${digits.slice(8,10)}`;
        }
        return phone;
    },

    // Setup input mask for phone fields
    setupMask(input) {
        if (!input) return;

        input.addEventListener('input', function(e) {
            let value = e.target.value.replace(/\D/g, '');

            if (value.length === 0) {
                e.target.value = '';
                return;
            }

            // Start with +7 if entering from scratch
            if (value[0] !== '7' && value.length === 1) {
                value = '7' + value;
            }

            let formatted = '+7';
            if (value.length > 1) {
                formatted += ' (' + value.slice(1, 4);
            }
            if (value.length > 4) {
                formatted += ') ' + value.slice(4, 7);
            }
            if (value.length > 7) {
                formatted += '-' + value.slice(7, 9);
            }
            if (value.length > 9) {
                formatted += '-' + value.slice(9, 11);
            }

            e.target.value = formatted;
        });

        // Set placeholder
        input.placeholder = '+7 (___) ___-__-__';
    }
};

// API utilities
const API = {
    async get(url) {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    },

    async post(url, data) {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    },

    async put(url, data) {
        const response = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    },

    async delete(url) {
        const response = await fetch(url, { method: 'DELETE' });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    }
};

// Loading state utilities
const Loading = {
    show(container, message = 'Загрузка...') {
        if (!container) return;
        container.innerHTML = `
            <div class="loading">
                <div class="loading-spinner"></div>
                <div>${message}</div>
            </div>
        `;
    },

    showInTable(tbody, colspan = 7) {
        if (!tbody) return;
        tbody.innerHTML = `
            <tr>
                <td colspan="${colspan}">
                    <div class="loading">
                        <div class="loading-spinner"></div>
                        <div>Загрузка данных...</div>
                    </div>
                </td>
            </tr>
        `;
    }
};

// Empty state utilities
const EmptyState = {
    show(container, icon = '📭', message = 'Нет данных') {
        if (!container) return;
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">${icon}</div>
                <div class="empty-state-text">${message}</div>
            </div>
        `;
    },

    showInTable(tbody, colspan = 7, icon = '📭', message = 'Нет записей') {
        if (!tbody) return;
        tbody.innerHTML = `
            <tr>
                <td colspan="${colspan}">
                    <div class="empty-state">
                        <div class="empty-state-icon">${icon}</div>
                        <div class="empty-state-text">${message}</div>
                    </div>
                </td>
            </tr>
        `;
    }
};

// Modal utilities
const Modal = {
    open(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'block';
            document.body.style.overflow = 'hidden';
        }
    },

    close(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'none';
            document.body.style.overflow = '';
        }
    },

    // Setup click outside to close
    setupClickOutside(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    this.close(modalId);
                }
            });
        }
    }
};

// Status badge renderer
function renderStatusBadge(status) {
    const statusMap = {
        'scheduled': { text: 'Запланировано', class: 'status-scheduled' },
        'completed': { text: 'Завершено', class: 'status-completed' },
        'cancelled': { text: 'Отменено', class: 'status-cancelled' }
    };

    const info = statusMap[status] || { text: status, class: '' };
    return `<span class="status-badge ${info.class}">${info.text}</span>`;
}

// Initialize common functionality
document.addEventListener('DOMContentLoaded', () => {
    // Setup phone masks on all phone inputs
    document.querySelectorAll('input[type="tel"]').forEach(input => {
        Phone.setupMask(input);
    });

    // Setup modal click outside
    document.querySelectorAll('.modal').forEach(modal => {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.style.display = 'none';
                document.body.style.overflow = '';
            }
        });
    });
});
