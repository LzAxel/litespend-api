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
  description: string;
  amount: string;
  date: string;
  created_at: string;
}

export interface CreateTransactionRequest {
  category_id: number;
  description: string;
  amount: string;
  date: string;
}

export interface UpdateTransactionRequest {
  category_id?: number;
  description?: string;
  amount?: string;
  date?: string;
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
  free_balance: string;
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

export type SortField = 'date' | 'description' | 'category';
export type SortOrder = 'asc' | 'desc';

export const transactionsApi = {
  getAll: (page = 1, limit = 10, sortBy?: SortField, sortOrder?: SortOrder, search?: string) => {
    const params: Record<string, string> = { page: String(page), limit: String(limit) };
    if (sortBy) params.sort_by = sortBy;
    if (sortOrder) params.sort_order = sortOrder;
    if (search) params.search = search;
    return api.get<PaginatedResponse<Transaction>>('/transactions', { params });
  },
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

export type FrequencyType = 0 | 1 | 2 | 3;

export interface PrescribedExpanse {
  id: number;
  user_id: number;
  category_id: number;
  description: string;
  frequency: FrequencyType;
  amount: string;
  date: string;
  created_at: string;
}

export interface PrescribedExpanseWithPaymentStatus extends PrescribedExpanse {
  is_paid: boolean;
  paid_amount: string;
  transaction_id?: number;
}

export interface CreatePrescribedExpanseRequest {
  category_id: number;
  description: string;
  frequency: FrequencyType;
  amount: string;
  date: string;
}

export interface UpdatePrescribedExpanseRequest {
  category_id?: number;
  description?: string;
  frequency?: FrequencyType;
  amount?: string;
  date?: string;
}

export const prescribedExpansesApi = {
  getAll: () => api.get<PrescribedExpanse[]>('/prescribed-expanses'),
  getAllWithPaymentStatus: () =>
    api.get<PrescribedExpanseWithPaymentStatus[]>('/prescribed-expanses/with-payment-status'),
  getById: (id: number) => api.get<PrescribedExpanse>(`/prescribed-expanses/${id}`),
  create: (data: CreatePrescribedExpanseRequest) =>
    api.post<{ id: number }>('/prescribed-expanses', data),
  update: (id: number, data: UpdatePrescribedExpanseRequest) =>
    api.put(`/prescribed-expanses/${id}`, data),
  delete: (id: number) => api.delete(`/prescribed-expanses/${id}`),
  markAsPaid: (id: number) =>
    api.post<{ transaction_id: number; message: string }>(`/prescribed-expanses/${id}/mark-as-paid`),
  markAsPaidPartial: (id: number, amount: string) =>
    api.post<{ transaction_id: number; message: string }>(`/prescribed-expanses/${id}/mark-as-paid-partial`, {
      amount,
    }),
};

export interface ExcelColumnMapping {
  transaction_description?: string;
  transaction_amount?: string;
  transaction_type?: string;
  transaction_date?: string;
  transaction_category?: string;
  category_name?: string;
  category_type?: string;
  prescribed_expanse_description?: string;
  prescribed_expanse_amount?: string;
  prescribed_expanse_frequency?: string;
  prescribed_expanse_date?: string;
  prescribed_expanse_category?: string;
}

export interface ExcelFileStructure {
  columns: string[];
  rows: number;
}

export interface ImportResult {
  transactions_created: number;
  categories_created: number;
  prescribed_expanses_created: number;
  errors?: string[];
}

export const importApi = {
  parseFile: (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post<ExcelFileStructure>('/import/parse', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  importData: (file: File, mapping: ExcelColumnMapping) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('mapping', JSON.stringify(mapping));
    return api.post<ImportResult>('/import/data', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};

