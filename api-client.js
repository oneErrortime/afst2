/**
 * API Client - централизованное управление API запросами
 */

class ApiClient {
    constructor() {
        this.baseURL = localStorage.getItem('apiUrl') || '';
        this.token = localStorage.getItem('authToken') || null;
    }

    /**
     * Устанавливает базовый URL API
     * @param {string} url
     */
    setBaseURL(url) {
        try {
            this.baseURL = Sanitizer.sanitizeUrl(url);
            localStorage.setItem('apiUrl', this.baseURL);
        } catch (error) {
            throw new Error('Invalid API URL');
        }
    }

    /**
     * Получает текущий базовый URL
     * @returns {string}
     */
    getBaseURL() {
        return this.baseURL;
    }

    /**
     * Устанавливает токен аутентификации
     * @param {string} token
     */
    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('authToken', token);
        } else {
            localStorage.removeItem('authToken');
        }
    }

    /**
     * Получает токен
     * @returns {string|null}
     */
    getToken() {
        return this.token;
    }

    /**
     * Очищает токен
     */
    clearToken() {
        this.token = null;
        localStorage.removeItem('authToken');
    }

    /**
     * Проверяет, аутентифицирован ли пользователь
     * @returns {boolean}
     */
    isAuthenticated() {
        return !!this.token;
    }

    /**
     * Выполняет HTTP запрос
     * @private
     */
    async request(endpoint, options = {}) {
        if (!this.baseURL) {
            throw new Error('API URL not configured');
        }

        const url = `${this.baseURL}${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        // Добавляем токен если есть
        if (this.token && !options.skipAuth) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                // Проверяем, является ли data объектом или строкой
                let errorMessage = 'Request failed';
                
                if (typeof data === 'object') {
                    // Пытаемся получить detail из объекта
                    if (data.detail) {
                        errorMessage = typeof data.detail === 'string' 
                            ? data.detail 
                            : (Array.isArray(data.detail) 
                                ? data.detail.map(err => err.msg || err).join(', ') 
                                : JSON.stringify(data.detail));
                    } else if (data.message) {
                        errorMessage = data.message;
                    } else {
                        // Если detail нет, пробуем получить сообщения из полей ошибок валидации
                        const fieldErrors = Object.keys(data)
                            .filter(key => key !== 'status_code')
                            .map(key => `${key}: ${Array.isArray(data[key]) ? data[key].join(', ') : data[key]}`)
                            .join('; ');
                        
                        if (fieldErrors) {
                            errorMessage = fieldErrors;
                        }
                    }
                } else if (typeof data === 'string') {
                    errorMessage = data;
                }
                
                throw new ApiError(
                    errorMessage,
                    response.status,
                    data
                );
            }

            return data;
        } catch (error) {
            if (error instanceof ApiError) {
                throw error;
            }
            
            // Ошибки сети
            if (error.name === 'TypeError') {
                throw new ApiError('Network error. Please check your connection.', 0);
            }

            throw new ApiError(error.message || 'Unknown error occurred', 0);
        }
    }

    /**
     * GET запрос
     */
    async get(endpoint, options = {}) {
        return this.request(endpoint, {
            method: 'GET',
            ...options
        });
    }

    /**
     * POST запрос
     */
    async post(endpoint, data, options = {}) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
            ...options
        });
    }

    /**
     * PUT запрос
     */
    async put(endpoint, data, options = {}) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
            ...options
        });
    }

    /**
     * DELETE запрос
     */
    async delete(endpoint, options = {}) {
        return this.request(endpoint, {
            method: 'DELETE',
            ...options
        });
    }

    // ========== AUTH API ==========

    /**
     * Регистрация пользователя
     * @param {UserCreateDTO} userDTO
     * @returns {Promise<UserDTO>}
     */
    async register(userDTO) {
        userDTO.validate();
        const data = await this.post('/auth/register', userDTO.toJSON(), { skipAuth: true });
        return new DTO.UserDTO(data);
    }

    /**
     * Вход пользователя
     * @param {string} email
     * @param {string} password
     * @returns {Promise<TokenResponseDTO>}
     */
    async login(email, password) {
        // Отправляем как form data для совместимости с OAuth2PasswordRequestForm
        const formData = new URLSearchParams();
        formData.append('username', email);  // OAuth2 использует 'username' для email
        formData.append('password', password);
        
        const response = await fetch(`${this.baseURL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            // Проверяем, является ли data объектом или строкой
            let errorMessage = 'Login failed';
            
            if (typeof data === 'object') {
                // Пытаемся получить detail из объекта
                if (data.detail) {
                    errorMessage = typeof data.detail === 'string' 
                        ? data.detail 
                        : (Array.isArray(data.detail) 
                            ? data.detail.map(err => err.msg || err).join(', ') 
                            : JSON.stringify(data.detail));
                } else if (data.message) {
                    errorMessage = data.message;
                } else {
                    // Если detail нет, пробуем получить сообщения из полей ошибок валидации
                    const fieldErrors = Object.keys(data)
                        .filter(key => key !== 'status_code')
                        .map(key => `${key}: ${Array.isArray(data[key]) ? data[key].join(', ') : data[key]}`)
                        .join('; ');
                    
                    if (fieldErrors) {
                        errorMessage = fieldErrors;
                    }
                }
            } else if (typeof data === 'string') {
                errorMessage = data;
            }
            
            throw new ApiError(
                errorMessage,
                response.status,
                data
            );
        }
        
        const tokenDTO = new DTO.TokenResponseDTO(data);
        tokenDTO.validate();
        
        // Сохраняем токен
        this.setToken(tokenDTO.access_token);
        
        return tokenDTO;
    }

    /**
     * Выход пользователя
     */
    logout() {
        this.clearToken();
    }

    // ========== BOOKS API ==========

    /**
     * Получить все книги
     * @returns {Promise<BookDTO[]>}
     */
    async getBooks() {
        const data = await this.get('/books/');
        return data.map(book => new DTO.BookDTO(book));
    }

    /**
     * Получить книгу по ID
     * @param {number} id
     * @returns {Promise<BookDTO>}
     */
    async getBook(id) {
        const data = await this.get(`/books/${id}`);
        return new DTO.BookDTO(data);
    }

    /**
     * Создать книгу
     * @param {BookDTO} bookDTO
     * @returns {Promise<BookDTO>}
     */
    async createBook(bookDTO) {
        bookDTO.validate();
        const data = await this.post('/books/', bookDTO.toJSON());
        return new DTO.BookDTO(data);
    }

    /**
     * Обновить книгу
     * @param {number} id
     * @param {BookDTO} bookDTO
     * @returns {Promise<BookDTO>}
     */
    async updateBook(id, bookDTO) {
        bookDTO.validate();
        const data = await this.put(`/books/${id}`, bookDTO.toJSON());
        return new DTO.BookDTO(data);
    }

    /**
     * Удалить книгу
     * @param {number} id
     * @returns {Promise<void>}
     */
    async deleteBook(id) {
        await this.delete(`/books/${id}`);
    }

    // ========== READERS API ==========

    /**
     * Получить всех читателей
     * @returns {Promise<ReaderDTO[]>}
     */
    async getReaders() {
        const data = await this.get('/readers/');
        return data.map(reader => new DTO.ReaderDTO(reader));
    }

    /**
     * Получить читателя по ID
     * @param {number} id
     * @returns {Promise<ReaderDTO>}
     */
    async getReader(id) {
        const data = await this.get(`/readers/${id}`);
        return new DTO.ReaderDTO(data);
    }

    /**
     * Создать читателя
     * @param {ReaderDTO} readerDTO
     * @returns {Promise<ReaderDTO>}
     */
    async createReader(readerDTO) {
        readerDTO.validate();
        const data = await this.post('/readers/', readerDTO.toJSON());
        return new DTO.ReaderDTO(data);
    }

    /**
     * Обновить читателя
     * @param {number} id
     * @param {ReaderDTO} readerDTO
     * @returns {Promise<ReaderDTO>}
     */
    async updateReader(id, readerDTO) {
        readerDTO.validate();
        const data = await this.put(`/readers/${id}`, readerDTO.toJSON());
        return new DTO.ReaderDTO(data);
    }

    /**
     * Удалить читателя
     * @param {number} id
     * @returns {Promise<void>}
     */
    async deleteReader(id) {
        await this.delete(`/readers/${id}`);
    }

    // ========== BORROWS API ==========

    /**
     * Выдать книгу
     * @param {BorrowDTO} borrowDTO
     * @returns {Promise<Object>}
     */
    async borrowBook(borrowDTO) {
        borrowDTO.validate();
        return await this.post('/borrows/borrow', borrowDTO.toJSON());
    }

    /**
     * Вернуть книгу
     * @param {number} bookId
     * @param {number} readerId
     * @returns {Promise<Object>}
     */
    async returnBook(bookId, readerId) {
        const borrowDTO = new DTO.BorrowDTO({ book_id: bookId, reader_id: readerId });
        borrowDTO.validate();
        return await this.post('/borrows/return', borrowDTO.toJSON());
    }

    /**
     * Получить выданные книги читателя
     * @param {number} readerId
     * @returns {Promise<BorrowDTO[]>}
     */
    async getReaderBorrows(readerId) {
        const data = await this.get(`/borrows/reader/${readerId}/borrowed`);
        return data.map(borrow => new DTO.BorrowDTO(borrow));
    }

    /**
     * Получить все записи о выдаче
     * @returns {Promise<BorrowDTO[]>}
     */
    async getAllBorrows() {
        const data = await this.get('/borrows/');
        return data.map(borrow => new DTO.BorrowDTO(borrow));
    }
}

/**
 * Кастомный класс ошибок API
 */
class ApiError extends Error {
    constructor(message, status, data = null) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
        this.data = data;
    }

    /**
     * Проверяет, является ли ошибка ошибкой аутентификации
     */
    isAuthError() {
        return this.status === 401 || this.status === 403;
    }

    /**
     * Проверяет, является ли ошибка ошибкой валидации
     */
    isValidationError() {
        return this.status === 422 || this.status === 400;
    }

    /**
     * Проверяет, является ли ошибка серверной ошибкой
     */
    isServerError() {
        return this.status >= 500;
    }
}

// Экспорт
if (typeof window !== 'undefined') {
    window.ApiClient = ApiClient;
    window.ApiError = ApiError;
}