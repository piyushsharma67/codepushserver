import axios from 'axios';

const API_URL = 'http://localhost:8080';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests if it exists
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export interface RegisterRequest {
    username: string;
    email: string;
    password: string;
    companyName?: string;
    phoneNumber?: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface User {
    id: number;
    username: string;
    email: string;
    companyName?: string;
    phoneNumber?: string;
    appId: string;
    token: string;
    createdAt: string;
}

export interface LoginResponse {
    token: string;
    expiresAt: string;
    user: User;
}

export const authService = {
    register: async (data: RegisterRequest) => {
        const response = await api.post<{ user: User }>('/v2/auth/register', data);
        return response.data;
    },

    login: async (data: LoginRequest) => {
        const response = await api.post<LoginResponse>('/v2/auth/login', data);
        return response.data;
    },
};

export const userService = {
    getProfile: async () => {
        const response = await api.get<User>('/v2/user/profile');
        return response.data;
    },

    updateProfile: async (data: Partial<User>) => {
        const response = await api.put<User>('/v2/user/profile', data);
        return response.data;
    },
}; 