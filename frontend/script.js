const API_URL = "http://localhost:8080";
const FALLBACK_IMAGE = `${API_URL}/static/images/cinema-night.svg`;
const SEATS = ["A", "B", "C", "D"].flatMap((row) => Array.from({ length: 8 }, (_, index) => `${row}${index + 1}`));

const moviesGrid = document.querySelector("#moviesGrid");
const moviesCount = document.querySelector("#moviesCount");
const statusBox = document.querySelector("#status");
const refreshButton = document.querySelector("#refreshButton");
const bookingDialog = document.querySelector("#bookingDialog");
const dialogTitle = document.querySelector("#dialogTitle");
const dialogSubtitle = document.querySelector("#dialogSubtitle");
const seatsGrid = document.querySelector("#seatsGrid");
const selectedSeatsText = document.querySelector("#selectedSeatsText");
const confirmBookingButton = document.querySelector("#confirmBookingButton");

let movies = [];
let activeMovie = null;
let activeBooking = null;
let selectedSeats = new Set();

function getBookings() {
    return JSON.parse(localStorage.getItem("movie-bookings") || "{}");
}

function saveBookings(bookings) {
    localStorage.setItem("movie-bookings", JSON.stringify(bookings));
}

function showStatus(message) {
    statusBox.textContent = message;
    statusBox.hidden = !message;
}

function getImageUrl(movie) {
    if (!movie.image_url) {
        return FALLBACK_IMAGE;
    }
    if (movie.image_url.startsWith("http")) {
        return movie.image_url;
    }
    return `${API_URL}${movie.image_url}`;
}

function formatDuration(seconds) {
    const minutes = Math.round(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const rest = minutes % 60;
    return hours > 0 ? `${hours} ч ${rest} мин` : `${minutes} мин`;
}

function formatDate(value) {
    if (!value) {
        return "";
    }
    return new Intl.DateTimeFormat("ru-RU", {
        day: "2-digit",
        month: "long",
        year: "numeric"
    }).format(new Date(value));
}

function getBookedSeats(movieId, ignoredBookingId = null) {
    const booking = getBookings()[movieId];
    if (booking?.id === ignoredBookingId) {
        return new Set();
    }
    return new Set(booking?.seats || []);
}

function renderMovies() {
    moviesGrid.innerHTML = "";
    moviesCount.textContent = `${movies.length} фильмов`;

    if (movies.length === 0) {
        showStatus("Фильмы не найдены");
        return;
    }

    showStatus("");
    const bookings = getBookings();

    movies.forEach((movie) => {
        const booking = bookings[movie.id];
        const card = document.createElement("article");
        card.className = "movie-card";

        const genres = (movie.genres || []).map((genre) => `<span class="genre">${genre.name}</span>`).join("");
        const bookingNote = booking ? `<p class="booking-note">Забронировано: ${booking.seats.join(", ")}</p>` : "";
        const actionButtons = booking
            ? `
                <button class="primary-button" data-action="edit" data-id="${movie.id}" type="button">Изменить места</button>
                <button class="danger-button" data-action="cancel" data-id="${movie.id}" type="button">Отменить бронь</button>
            `
            : `<button class="primary-button" data-action="book" data-id="${movie.id}" type="button">Забронировать</button>`;

        card.innerHTML = `
            <img class="poster" src="${getImageUrl(movie)}" alt="${movie.title}">
            <div class="movie-content">
                <div class="movie-title-row">
                    <h2 class="movie-title">${movie.title}</h2>
                    <span class="age">${movie.age_rating}+</span>
                </div>
                <p class="meta">${movie.director} · ${formatDuration(movie.duration)} · ${formatDate(movie.release_date)}</p>
                <p class="description">${movie.description || ""}</p>
                <div class="genres">${genres}</div>
                ${bookingNote}
                <div class="card-actions">${actionButtons}</div>
            </div>
        `;

        moviesGrid.append(card);
    });
}

function renderSeats() {
    seatsGrid.innerHTML = "";
    selectedSeatsText.textContent = selectedSeats.size > 0
        ? `Выбрано: ${Array.from(selectedSeats).join(", ")}`
        : "Места не выбраны";
    confirmBookingButton.disabled = selectedSeats.size === 0;

    const bookedSeats = getBookedSeats(activeMovie.id, activeBooking?.id);

    SEATS.forEach((seat) => {
        const button = document.createElement("button");
        button.type = "button";
        button.textContent = seat;
        button.className = "seat-button";

        if (bookedSeats.has(seat)) {
            button.classList.add("booked");
            button.disabled = true;
        }
        if (selectedSeats.has(seat)) {
            button.classList.add("selected");
        }

        button.addEventListener("click", () => {
            if (selectedSeats.has(seat)) {
                selectedSeats.delete(seat);
            } else {
                selectedSeats.add(seat);
            }
            renderSeats();
        });

        seatsGrid.append(button);
    });
}

function openBooking(movieId) {
    activeMovie = movies.find((movie) => movie.id === movieId);
    activeBooking = getBookings()[movieId] || null;
    selectedSeats = new Set(activeBooking?.seats || []);
    confirmBookingButton.textContent = activeBooking ? "Сохранить изменения" : "Забронировать";
    dialogTitle.textContent = activeMovie.title;
    dialogSubtitle.textContent = `${formatDuration(activeMovie.duration)} · ${formatDate(activeMovie.release_date)}`;
    renderSeats();
    bookingDialog.showModal();
}

async function saveBooking() {
    if (!activeMovie || selectedSeats.size === 0) {
        return;
    }

    confirmBookingButton.disabled = true;
    const seats = Array.from(selectedSeats);
    const isUpdate = Boolean(activeBooking);
    const url = isUpdate
        ? `${API_URL}/api/v1/booking/${activeBooking.id}`
        : `${API_URL}/api/v1/booking`;
    const payload = isUpdate
        ? { seats }
        : { movie_id: activeMovie.id, seats };

    const response = await fetch(url, {
        method: isUpdate ? "PUT" : "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    });

    if (!response.ok) {
        confirmBookingButton.disabled = false;
        showStatus(isUpdate ? "Не удалось изменить бронь" : "Не удалось создать бронь");
        return;
    }

    const booking = await response.json();
    const bookings = getBookings();
    bookings[activeMovie.id] = {
        id: booking.id,
        seats: booking.seats
    };
    saveBookings(bookings);

    bookingDialog.close();
    renderMovies();
}

async function cancelBooking(movieId) {
    const bookings = getBookings();
    const booking = bookings[movieId];
    if (!booking) {
        return;
    }

    const response = await fetch(`${API_URL}/api/v1/booking/${booking.id}`, {
        method: "DELETE"
    });

    if (!response.ok && response.status !== 404) {
        showStatus("Не удалось отменить бронь");
        return;
    }

    delete bookings[movieId];
    saveBookings(bookings);
    renderMovies();
}

async function loadMovies() {
    showStatus("Загрузка фильмов");
    moviesGrid.innerHTML = "";

    try {
        const response = await fetch(`${API_URL}/api/v1/movies`);
        if (!response.ok) {
            throw new Error("movies request failed");
        }
        movies = await response.json();
        renderMovies();
    } catch {
        movies = [];
        moviesCount.textContent = "Фильмы";
        showStatus("Backend недоступен");
    }
}

moviesGrid.addEventListener("click", (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) {
        return;
    }

    const movieId = Number(button.dataset.id);
    if (button.dataset.action === "book") {
        openBooking(movieId);
    }
    if (button.dataset.action === "edit") {
        openBooking(movieId);
    }
    if (button.dataset.action === "cancel") {
        cancelBooking(movieId);
    }
});

confirmBookingButton.addEventListener("click", saveBooking);
refreshButton.addEventListener("click", loadMovies);

loadMovies();
