import axios from 'axios';

// Temporary simple types to test import
interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

interface User {
  id: string;
  email: string;
  name: string;
  picture?: string;
  company_id: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  onboarded: boolean;
  invitation_status: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

class ApiService {
  private baseURL = 'http://localhost:8080/api/v1';

  private getHeaders() {
    const token = localStorage.getItem('auth_token');
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  // Auth endpoints
  async login(data: LoginRequest): Promise<APIResponse<{ token: string; user: User }>> {
    const response = await axios.post(`${this.baseURL}/auth/login`, data, {
      headers: this.getHeaders(),
    });
    return response.data;
  }

  async register(data: RegisterRequest): Promise<APIResponse<{ token: string; user: User }>> {
    const response = await axios.post(`${this.baseURL}/auth/register`, data, {
      headers: this.getHeaders(),
    });
    return response.data;
  }

  async verifyToken(): Promise<APIResponse<{ user: User }>> {
    const response = await axios.get(`${this.baseURL}/auth/verify`, {
      headers: this.getHeaders(),
    });
    return response.data;
  }

  async logout(): Promise<APIResponse> {
    const response = await axios.post(`${this.baseURL}/auth/logout`, {}, {
      headers: this.getHeaders(),
    });
    return response.data;
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string; service: string }> {
    const response = await axios.get(`${this.baseURL}/health`);
    return response.data;
  }
}

export const apiService = new ApiService();
export default apiService;
