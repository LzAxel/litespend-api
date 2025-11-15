import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8888/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export type TransactionType = 0 | 1;

export interface Transaction {
  id: number;
  user_id: number;
  category_id: number;
  goal_id?: number;
  description: string;
  amount: string;
  type: TransactionType;
  date_time: string;
  created_at: string;
}

export interface CreateTransactionRequest {
  category_id: number;
  goal_id?: number;
  description: string;
  amount: string;
  type: TransactionType;
  date_time: string;
}

export interface UpdateTransactionRequest {
  category_id?: number;
  goal_id?: number;
  description?: string;
  amount?: string;
  type?: TransactionType;
  date_time?: string;
}

export interface Category {
  id: number;
  user_id: number;
  name: string;
  type: TransactionType;
  created_at: string;
}

export interface CreateCategoryRequest {
  name: string;
  type: TransactionType;
}

export interface UpdateCategoryRequest {
  name?: string;
  type?: TransactionType;
}

export interface User {
  id: number;
  username: string;
  role: number;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface CurrentBalanceStatistics {
  balance: string;
  total_income: string;
  total_expense: string;
}

export type PeriodType = 'day' | 'week' | 'month';

export interface PeriodStatisticsItem {
  period: string;
  income: string;
  expense: string;
  balance: string;
}

export interface PeriodStatisticsResponse {
  period: PeriodType;
  items: PeriodStatisticsItem[];
}

export interface CategoryStatisticsItem {
  category_id: number;
  category_name: string;
  period: string;
  income: string;
  expense: string;
}

export interface CategoryStatisticsResponse {
  period: PeriodType;
  items: CategoryStatisticsItem[];
}

export const authApi = {
  login: (data: LoginRequest) => api.post('/user/login', data),
  register: (data: RegisterRequest) => api.post('/user/register', data),
  logout: () => api.post('/user/logout'),
};

export const transactionsApi = {
  getAll: (page = 1, limit = 10) =>
    api.get<PaginatedResponse<Transaction>>('/transactions', { params: { page, limit } }),
  getById: (id: number) => api.get<Transaction>(`/transactions/${id}`),
  create: (data: CreateTransactionRequest) => api.post<{ id: number }>('/transactions', data),
  update: (id: number, data: UpdateTransactionRequest) =>
    api.put(`/transactions/${id}`, data),
  delete: (id: number) => api.delete(`/transactions/${id}`),
  getBalanceStatistics: () => api.get<CurrentBalanceStatistics>('/transactions/statistics/balance'),
  getPeriodStatistics: (period: PeriodType, from?: string, to?: string) =>
    api.get<PeriodStatisticsResponse>('/transactions/statistics/periods', {
      params: { period, from, to },
    }),
  getCategoryStatistics: (period: PeriodType, from?: string, to?: string) =>
    api.get<CategoryStatisticsResponse>('/transactions/statistics/categories', {
      params: { period, from, to },
    }),
};

export const categoriesApi = {
  getAll: () => api.get<Category[]>('/categories'),
  getById: (id: number) => api.get<Category>(`/categories/${id}`),
  create: (data: CreateCategoryRequest) => api.post<{ id: number }>('/categories', data),
  update: (id: number, data: UpdateCategoryRequest) => api.put(`/categories/${id}`, data),
  delete: (id: number) => api.delete(`/categories/${id}`),
  getByType: (type: TransactionType) =>
    api.get<Category[]>(`/categories?type=${type}`),
};

