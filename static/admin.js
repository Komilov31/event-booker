// API Configuration
const API_BASE = window.location.origin;

// DOM Elements
const createEventForm = document.getElementById('createEventForm');
const eventsList = document.getElementById('eventsList');
const refreshEventsBtn = document.getElementById('refreshEvents');
const eventModal = document.getElementById('eventModal');
const modalTitle = document.getElementById('modalTitle');
const modalBody = document.getElementById('modalBody');
const closeModal = document.querySelector('.close');
const refreshBookingsBtn = document.getElementById('refreshBookings');
const statusFilter = document.getElementById('statusFilter');
const bulkConfirmPaymentBtn = document.getElementById('bulkConfirmPayment');
const selectAllCheckbox = document.getElementById('selectAll');

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    loadEvents();
    loadAllBookings();
});
refreshEventsBtn.addEventListener('click', loadEvents);
refreshBookingsBtn.addEventListener('click', loadAllBookings);
createEventForm.addEventListener('submit', handleCreateEvent);
closeModal.addEventListener('click', () => eventModal.style.display = 'none');
statusFilter.addEventListener('change', filterBookings);
bulkConfirmPaymentBtn.addEventListener('click', bulkConfirmPayment);
selectAllCheckbox.addEventListener('change', toggleAllCheckboxes);

// Functions
async function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 5000);
}

async function loadEvents() {
    try {
        eventsList.innerHTML = '<div class="loading">Загрузка мероприятий...</div>';

        const response = await fetch(`${API_BASE}/events`);
        if (!response.ok) throw new Error('Failed to load events');

        const events = await response.json();
        displayEvents(events);
    } catch (error) {
        console.error('Error loading events:', error);
        eventsList.innerHTML = '<div class="error">Ошибка загрузки мероприятий</div>';
        showNotification('Ошибка загрузки мероприятий', 'error');
    }
}

function displayEvents(events) {
    if (events.length === 0) {
        eventsList.innerHTML = '<div class="no-events">Нет доступных мероприятий</div>';
        return;
    }

    eventsList.innerHTML = events.map(event => createEventCard(event)).join('');
}

function createEventCard(event) {
    const availableSeats = event.all_seats - event.booked;
    const isAvailable = availableSeats > 0;
    const bookingsCount = event.bookings ? event.bookings.length : 0;

    return `
        <div class="event-card ${!isAvailable ? 'unavailable' : ''}">
            <div class="event-header">
                <h3 class="event-title">${event.event_name}</h3>
                <span class="event-id">#${event.id}</span>
            </div>

            <div class="event-meta">
                <span>📅 Создано: ${new Date(event.created_at).toLocaleDateString('ru-RU')}</span>
                <span>📋 Бронирований: ${bookingsCount}</span>
            </div>

            <div class="event-stats">
                <div class="stat">
                    <div class="stat-value">${event.all_seats}</div>
                    <div class="stat-label">Всего мест</div>
                </div>
                <div class="stat">
                    <div class="stat-value">${event.booked}</div>
                    <div class="stat-label">Забронировано</div>
                </div>
                <div class="stat">
                    <div class="stat-value" style="color: ${availableSeats > 0 ? '#2ecc71' : '#e74c3c'}">${availableSeats}</div>
                    <div class="stat-label">Доступно</div>
                </div>
            </div>

            <div class="event-actions">
                <button onclick="showEventDetails(${event.id})" class="btn btn-secondary">
                    📊 Подробности и бронирования
                </button>
                ${!isAvailable ? '<span class="no-seats">Мест нет</span>' : ''}
            </div>
        </div>
    `;
}

async function showEventDetails(eventId) {
    try {
        const response = await fetch(`${API_BASE}/events/${eventId}`);
        if (!response.ok) throw new Error('Failed to load event details');

        const event = await response.json();

        modalTitle.textContent = event.event_name;
        modalBody.innerHTML = createEventDetailsHTML(event);

        eventModal.style.display = 'block';
    } catch (error) {
        console.error('Error loading event details:', error);
        showNotification('Ошибка загрузки деталей мероприятия', 'error');
    }
}

function createEventDetailsHTML(event) {
    const availableSeats = event.all_seats - event.booked;

    let bookingsHTML = '';
    if (event.bookings && event.bookings.length > 0) {
        bookingsHTML = `
            <h4>📋 Все бронирования (${event.bookings.length})</h4>
            <div class="bookings-table">
                ${event.bookings.map(booking => `
                    <div class="booking-row">
                        <div class="booking-info">
                            <div class="booking-id">ID: ${booking.id}</div>
                            <div class="booking-seats">Мест: ${booking.places_count}</div>
                        </div>
                        <div class="booking-status">
                            <span class="status-${booking.status}">${getStatusText(booking.status)}</span>
                        </div>
                        <div class="booking-time">
                            ${new Date(booking.created_at).toLocaleString('ru-RU')}
                        </div>
                        <div class="booking-actions">
                            ${booking.status !== 'paid' ? `<button onclick="confirmPayment(${booking.id})" class="btn btn-success btn-small">✅ Подтвердить оплату</button>` : '<span class="payment-confirmed">✅ Оплата подтверждена</span>'}
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
    } else {
        bookingsHTML = '<div class="no-bookings">Бронирований пока нет</div>';
    }

    return `
        <div class="event-details">
            <div class="event-overview">
                <div class="stat-large">
                    <div class="stat-value">${event.all_seats}</div>
                    <div class="stat-label">Всего мест</div>
                </div>
                <div class="stat-large">
                    <div class="stat-value" style="color: ${event.booked > 0 ? '#3498db' : '#95a5a6'}">${event.booked}</div>
                    <div class="stat-label">Забронировано</div>
                </div>
                <div class="stat-large">
                    <div class="stat-value" style="color: ${availableSeats > 0 ? '#2ecc71' : '#e74c3c'}">${availableSeats}</div>
                    <div class="stat-label">Доступно</div>
                </div>
            </div>

            <div class="event-info">
                <p><strong>Название:</strong> ${event.event_name}</p>
                <p><strong>Создано:</strong> ${new Date(event.created_at).toLocaleString('ru-RU')}</p>
                <p><strong>Дата проведения:</strong> ${new Date(event.event_at).toLocaleString('ru-RU')}</p>
                <p><strong>Заполненность:</strong> ${Math.round((event.booked / event.all_seats) * 100)}%</p>
            </div>

            ${bookingsHTML}
        </div>

        <style>
            .event-details {
                max-height: 70vh;
                overflow-y: auto;
            }

            .event-overview {
                display: grid;
                grid-template-columns: repeat(3, 1fr);
                gap: 20px;
                margin-bottom: 30px;
                padding: 20px;
                background: #f8f9fa;
                border-radius: 10px;
            }

            .stat-large {
                text-align: center;
                padding: 20px;
                background: white;
                border-radius: 8px;
                box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            }

            .stat-large .stat-value {
                font-size: 2rem;
                font-weight: 700;
                margin-bottom: 5px;
            }

            .event-info {
                margin-bottom: 30px;
                padding: 20px;
                background: #e8f4f8;
                border-radius: 8px;
            }

            .event-info p {
                margin-bottom: 10px;
                font-size: 1.1rem;
            }

            .bookings-table {
                margin-top: 20px;
            }

            .booking-row {
                display: grid;
                grid-template-columns: 2fr 1fr 2fr 1fr;
                gap: 15px;
                padding: 15px;
                background: white;
                border-radius: 8px;
                margin-bottom: 10px;
                box-shadow: 0 2px 5px rgba(0,0,0,0.1);
                align-items: center;
            }

            .booking-info {
                display: flex;
                flex-direction: column;
                gap: 5px;
            }

            .booking-id {
                font-weight: 600;
                color: #2c3e50;
            }

            .booking-seats {
                color: #7f8c8d;
            }

            .booking-status {
                text-align: center;
            }

            .booking-time {
                color: #7f8c8d;
                font-size: 0.9rem;
            }

            .booking-actions {
                text-align: center;
            }

            .no-bookings {
                text-align: center;
                padding: 40px;
                color: #7f8c8d;
                font-style: italic;
                background: #f8f9fa;
                border-radius: 8px;
                margin: 20px 0;
            }

            .no-events {
                text-align: center;
                padding: 40px;
                color: #7f8c8d;
                font-style: italic;
                grid-column: 1 / -1;
            }

            .error {
                text-align: center;
                padding: 40px;
                color: #e74c3c;
                font-weight: 600;
                grid-column: 1 / -1;
            }

            .btn-small {
                padding: 5px 10px;
                font-size: 0.8rem;
                margin-left: 10px;
            }

            .payment-confirmed {
                color: #2ecc71;
                font-weight: 600;
                margin-left: 10px;
            }

            .status-pending {
                color: #f39c12;
                font-weight: 600;
            }

            .status-paid {
                color: #2ecc71;
                font-weight: 600;
            }

            .status-cancelled {
                color: #e74c3c;
                font-weight: 600;
            }
        </style>
    `;
}

function getStatusText(status) {
    const statusMap = {
        'pending': 'Ожидает оплаты',
        'paid': 'Оплачено',
        'cancelled': 'Отменено'
    };
    return statusMap[status] || status;
}

async function handleCreateEvent(e) {
    e.preventDefault();

    const formData = new FormData(e.target);
    const eventDateValue = formData.get('eventDate');

    if (!eventDateValue) {
        showNotification('Пожалуйста, выберите дату и время', 'error');
        return;
    }

    // Конвертируем datetime-local формат в ISO формат с часовым поясом
    const eventDate = new Date(eventDateValue);
    const isoDateString = eventDate.toISOString();

    const eventData = {
        event_name: formData.get('eventName'),
        event_at: isoDateString,
        all_seats: parseInt(formData.get('allSeats'))
    };

    if (!eventData.event_name || !eventData.event_at || !eventData.all_seats) {
        showNotification('Пожалуйста, заполните все поля', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(eventData)
        });

        if (!response.ok) throw new Error('Failed to create event');

        const newEvent = await response.json();

        showNotification(`Мероприятие "${newEvent.event_name}" успешно создано!`, 'success');
        createEventForm.reset();
        loadEvents();
    } catch (error) {
        console.error('Error creating event:', error);
        showNotification('Ошибка создания мероприятия', 'error');
    }
}

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === eventModal) {
        eventModal.style.display = 'none';
    }
});

// Payment confirmation function
async function confirmPayment(bookingId) {
    if (!confirm('Вы уверены, что хотите подтвердить оплату для этого бронирования?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events/${bookingId}/confirm`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) throw new Error('Failed to confirm payment');

        const result = await response.json();
        showNotification('Оплата успешно подтверждена!', 'success');

        // Close modal and reload events to show updated status
        eventModal.style.display = 'none';
        loadEvents();
    } catch (error) {
        console.error('Error confirming payment:', error);
        showNotification('Ошибка подтверждения оплаты', 'error');
    }
}

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && eventModal.style.display === 'block') {
        eventModal.style.display = 'none';
    }
});

// New functions for all bookings management
async function loadAllBookings() {
    try {
        const allBookingsTable = document.getElementById('allBookingsTable');
        allBookingsTable.innerHTML = '<div class="loading">Загрузка бронирований...</div>';

        const response = await fetch(`${API_BASE}/events`);
        if (!response.ok) throw new Error('Failed to load events');

        const events = await response.json();
        displayAllBookings(events);
    } catch (error) {
        console.error('Error loading all bookings:', error);
        document.getElementById('allBookingsTable').innerHTML = '<div class="error">Ошибка загрузки бронирований</div>';
        showNotification('Ошибка загрузки бронирований', 'error');
    }
}

function displayAllBookings(events) {
    const allBookings = [];

    events.forEach(event => {
        if (event.bookings && event.bookings.length > 0) {
            event.bookings.forEach(booking => {
                allBookings.push({
                    ...booking,
                    event_name: event.event_name,
                    event_id: event.id
                });
            });
        }
    });

    if (allBookings.length === 0) {
        document.getElementById('allBookingsTable').innerHTML = '<div class="no-bookings">Бронирований пока нет</div>';
        return;
    }

    const tableHTML = createAllBookingsTableHTML(allBookings);
    document.getElementById('allBookingsTable').innerHTML = tableHTML;
}

function createAllBookingsTableHTML(bookings) {
    return `
        <table class="bookings-table">
            <thead>
                <tr>
                    <th><input type="checkbox" id="selectAll"></th>
                    <th>ID</th>
                    <th>Мероприятие</th>
                    <th>Мест</th>
                    <th>Статус</th>
                    <th>Создано</th>
                    <th>Действия</th>
                </tr>
            </thead>
            <tbody>
                ${bookings.map(booking => `
                    <tr class="booking-row" data-status="${booking.status}">
                        <td><input type="checkbox" class="booking-checkbox" value="${booking.id}"></td>
                        <td>${booking.id}</td>
                        <td>${booking.event_name}</td>
                        <td>${booking.places_count}</td>
                        <td><span class="status-${booking.status}">${getStatusText(booking.status)}</span></td>
                        <td>${new Date(booking.created_at).toLocaleString('ru-RU')}</td>
                        <td>
                            ${booking.status !== 'paid' ?
                                `<button onclick="confirmPayment(${booking.id})" class="btn btn-success btn-small">✅ Подтвердить</button>` :
                                '<span class="payment-confirmed">✅ Подтверждено</span>'}
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

function filterBookings() {
    const filter = statusFilter.value;
    const rows = document.querySelectorAll('.booking-row');

    rows.forEach(row => {
        if (filter === 'all' || row.dataset.status === filter) {
            row.style.display = 'table-row';
        } else {
            row.style.display = 'none';
        }
    });
}

function toggleAllCheckboxes() {
    const checkboxes = document.querySelectorAll('.booking-checkbox');
    checkboxes.forEach(cb => cb.checked = selectAllCheckbox.checked);
}

async function bulkConfirmPayment() {
    const selectedBookings = document.querySelectorAll('.booking-checkbox:checked');
    if (selectedBookings.length === 0) {
        showNotification('Выберите бронирования для подтверждения', 'error');
        return;
    }

    if (!confirm(`Подтвердить оплату для ${selectedBookings.length} бронирований?`)) {
        return;
    }

    const bookingIds = Array.from(selectedBookings).map(cb => parseInt(cb.value));

    try {
        for (const bookingId of bookingIds) {
            const response = await fetch(`${API_BASE}/events/${bookingId}/confirm`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (!response.ok) throw new Error(`Failed to confirm payment for booking ${bookingId}`);
        }

        showNotification(`Успешно подтверждено ${bookingIds.length} бронирований!`, 'success');

        // Close modal and reload data to show updated status
        if (eventModal.style.display === 'block') {
            eventModal.style.display = 'none';
        }
        loadEvents();
        loadAllBookings();
    } catch (error) {
        console.error('Error bulk confirming payments:', error);
        showNotification('Ошибка при подтверждении оплат', 'error');
    }
}
