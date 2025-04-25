import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export interface LoginResponse {
    token: string;
    user: {
        id: string;
        email: string;
        username: string;
        companyName?: string;
        phoneNumber?: string;
    };
}

export const authService = {
    login: async (email: string, password: string): Promise<LoginResponse> => {
        try {
            const response = await axios.post<LoginResponse>(`${API_URL}/v2/auth/login`, {
                email,
                password,
            });
            return response.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(error.response?.data?.message || 'Login failed');
            }
            throw error;
        }
    },

    register: async (email: string, password: string, username: string): Promise<void> => {
        try {
            await axios.post(`${API_URL}/v2/auth/register`, {
                email,
                password,
                username,
            });
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(error.response?.data?.message || 'Registration failed');
            }
            throw error;
        }
    },
}; 