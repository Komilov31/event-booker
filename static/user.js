// API Configuration
const API_BASE = window.location.origin;

// DOM Elements
const eventsList = document.getElementById('eventsList');
const bookingsList = document.getElementById('bookingsList');
const refreshEventsBtn = document.getElementById('refreshEvents');
const userTelegramIdInput = document.getElementById('userTelegramId');
const bookingModal = document.getElementById('bookingModal');
const confirmModal = document.getElementById('confirmModal');
const bookingForm = document.getElementById('bookingForm');
const confirmPaymentBtn = document.getElementById('confirmPayment');
const closeModalBtns = document.querySelectorAll('.close, .close-modal');

// State
let currentBookingId = null;
let currentEventId = null;

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    loadEvents();
    loadUserBookings();
});

refreshEventsBtn.addEventListener('click', loadEvents);
bookingForm.addEventListener('submit', handleBooking);
confirmPaymentBtn.addEventListener('click', handlePaymentConfirmation);

closeModalBtns.forEach(btn => {
    btn.addEventListener('click', closeModals);
});

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

    return `
        <div class="event-card ${!isAvailable ? 'unavailable' : ''}">
            <div class="event-header">
                <h3 class="event-title">${event.event_name}</h3>
                <span class="event-id">#${event.id}</span>
            </div>

            <div class="event-meta">
                <span>📅 Создано: ${new Date(event.created_at).toLocaleDateString('ru-RU')}</span>
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
                ${isAvailable ?
                    `<button onclick="showBookingModal(${event.id}, '${event.event_name}', ${availableSeats})" class="btn btn-primary">
                        🎫 Забронировать
                    </button>` :
                    '<span class="no-seats">❌ Мест нет</span>'
                }
            </div>
        </div>
    `;
}

async function loadUserBookings() {
    const telegramId = userTelegramIdInput.value;
    if (!telegramId) {
        bookingsList.innerHTML = '<div class="no-bookings">Введите ваш Telegram ID для просмотра бронирований</div>';
        return;
    }

    try {
        bookingsList.innerHTML = '<div class="loading">Загрузка бронирований...</div>';

        const response = await fetch(`${API_BASE}/events`);
        if (!response.ok) throw new Error('Failed to load events');

        const events = await response.json();
        const userBookings = getUserBookingsFromEvents(events, telegramId);

        displayBookings(userBookings);
    } catch (error) {
        console.error('Error loading bookings:', error);
        bookingsList.innerHTML = '<div class="error">Ошибка загрузки бронирований</div>';
        showNotification('Ошибка загрузки бронирований', 'error');
    }
}

function getUserBookingsFromEvents(events, telegramId) {
    const bookings = [];

    events.forEach(event => {
        if (event.bookings) {
            event.bookings.forEach(booking => {
                if (booking.telegram_id == telegramId) {
                    bookings.push({
                        ...booking,
                        event_name: event.event_name
                    });
                }
            });
        }
    });

    return bookings.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
}

function displayBookings(bookings) {
    if (bookings.length === 0) {
        bookingsList.innerHTML = '<div class="no-bookings">У вас пока нет бронирований</div>';
        return;
    }

    bookingsList.innerHTML = bookings.map(booking => createBookingItem(booking)).join('');
}

function createBookingItem(booking) {
    const isUrgent = booking.status === 'pending' &&
                    (new Date() - new Date(booking.created_at)) > (10 * 60 * 1000); // 10 minutes

    return `
        <div class="booking-item ${isUrgent ? 'urgent' : ''}">
            <div class="booking-info">
                <div class="booking-title">${booking.event_name}</div>
                <div class="booking-details">
                    <span>ID брони: #${booking.id}</span> |
                    <span>Мест: ${booking.places_count}</span> |
                    <span>Создано: ${new Date(booking.created_at).toLocaleString('ru-RU')}</span>
                </div>
            </div>
            <div class="booking-actions">
                <span class="booking-status status-${booking.status}">${getStatusText(booking.status)}</span>
                ${booking.status === 'pending' ?
                    `<button onclick="showConfirmModal(${booking.id}, ${booking.event_id})" class="btn btn-success">
                        💳 Оплатить
                    </button>` :
                    ''
                }
            </div>
        </div>
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

function showBookingModal(eventId, eventName, availableSeats) {
    currentEventId = eventId;

    const modalTitle = document.getElementById('bookingModalTitle');
    const modalBody = document.getElementById('bookingModalBody');

    modalTitle.textContent = `Бронирование: ${eventName}`;
    modalBody.innerHTML = `
        <div class="form-group">
            <label for="placesCount">Количество мест:</label>
            <input type="number" id="placesCount" name="placesCount" min="1" max="${availableSeats}" required value="1">
            <small style="color: #7f8c8d;">Доступно: ${availableSeats} мест</small>
        </div>
        <div class="form-group">
            <label for="telegramId">Ваш Telegram ID:</label>
            <input type="number" id="telegramId" name="telegramId" required value="${userTelegramIdInput.value}">
        </div>
    `;

    bookingModal.style.display = 'block';
}

function showConfirmModal(bookingId, eventId) {
    currentBookingId = bookingId;
    currentEventId = eventId;

    const confirmModalTitle = document.getElementById('confirmModalTitle');
    const confirmModalBody = document.getElementById('confirmModalBody');

    confirmModalTitle.textContent = 'Подтверждение оплаты';
    confirmModalBody.innerHTML = `
        <p>Вы действительно хотите подтвердить оплату для бронирования #${bookingId}?</p>
        <p><strong>Внимание:</strong> После подтверждения оплата будет считаться завершенной.</p>
    `;

    confirmModal.style.display = 'block';
}

function closeModals() {
    bookingModal.style.display = 'none';
    confirmModal.style.display = 'none';
    currentBookingId = null;
    currentEventId = null;
}

async function handleBooking(e) {
    e.preventDefault();

    const formData = new FormData(e.target);
    const bookingData = {
        telegram_id: parseInt(formData.get('telegramId')),
        places_count: parseInt(formData.get('placesCount'))
    };

    if (!bookingData.telegram_id || !bookingData.places_count) {
        showNotification('Пожалуйста, заполните все поля', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/events/${currentEventId}/book`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookingData)
        });

        if (!response.ok) throw new Error('Failed to create booking');

        const booking = await response.json();

        showNotification(`Бронирование успешно создано! ID: #${booking.id}`, 'success');
        closeModals();
        loadEvents();
        loadUserBookings();
    } catch (error) {
        console.error('Error creating booking:', error);
        showNotification('Ошибка создания бронирования', 'error');
    }
}

async function handlePaymentConfirmation() {
    if (!currentBookingId) return;

    try {
        const response = await fetch(`${API_BASE}/events/${currentEventId}/confirm`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) throw new Error('Failed to confirm payment');

        showNotification('Оплата успешно подтверждена!', 'success');
        closeModals();
        loadEvents();
        loadUserBookings();
    } catch (error) {
        console.error('Error confirming payment:', error);
        showNotification('Ошибка подтверждения оплаты', 'error');
    }
}

// Auto-refresh bookings every 30 seconds
setInterval(() => {
    if (userTelegramIdInput.value) {
        loadUserBookings();
    }
}, 30000);

// Update bookings when Telegram ID changes
userTelegramIdInput.addEventListener('change', loadUserBookings);

// Close modals when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === bookingModal || e.target === confirmModal) {
        closeModals();
    }
});

// Close modals with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && (bookingModal.style.display === 'block' || confirmModal.style.display === 'block')) {
        closeModals();
    }
});
