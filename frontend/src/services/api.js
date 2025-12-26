import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If 401 and not already retried, try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const response = await axios.post(`${API_BASE_URL}/api/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token } = response.data;
        localStorage.setItem('access_token', access_token);

        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, logout user
        localStorage.clear();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  login: (username, password) =>
    api.post('/api/auth/login', { username, password }),
  
  getCurrentUser: () => api.get('/api/auth/me'),
  
  logout: () => {
    localStorage.clear();
  },
};

// Patient APIs
export const patientAPI = {
  getScans: () => api.get('/api/patient/scans'),
  
  getScanById: (id) => api.get(`/api/patient/scans/${id}`),
  
  getTreatmentPlans: () => api.get('/api/patient/plans'),
  
  getTreatmentPlan: (scanId) => api.get(`/api/patient/scans/${scanId}/plan`),
  
  getOffers: (planId) => api.get(`/api/patient/plans/${planId}/offers`),
  
  selectOffer: (offerId) => api.post('/api/patient/select-offer', { offer_id: offerId }),
  
  getAppointments: () => api.get('/api/patient/appointments'),
  
  createReview: (clinicId, rating, comment) =>
    api.post('/api/patient/reviews', { clinic_id: clinicId, rating, comment }),
  
  createComplaint: (clinicId, subject, description) =>
    api.post('/api/patient/complaints', { clinic_id: clinicId, subject, description }),
  
  updateSearchCriteria: (city, district, priceSegment) =>
    api.post('/api/patient/search-criteria', { city, district, price_segment: priceSegment }),
};

// Clinic APIs
export const clinicAPI = {
  getDashboard: (period = '30d') => api.get(`/api/clinic/dashboard?period=${period}`),
  
  getIncomingPlans: () => api.get('/api/clinic/incoming-plans'),
  
  createOffer: (offerData) => api.post('/api/clinic/offers', offerData),
  
  getLeads: () => api.get('/api/clinic/leads'),
  
  getAppointments: (status) =>
    api.get('/api/clinic/appointments', { params: { status } }),
  
  updateAppointment: (id, status, notes) =>
    api.put(`/api/clinic/appointments/${id}`, { status, notes }),
  
  getPriceList: (specialization) =>
    api.get('/api/clinic/price-list', { params: { specialization } }),
  
  updatePriceList: (items) => api.put('/api/clinic/price-list', items),
  
  getAnalytics: (period = '30d') => api.get(`/api/clinic/analytics?period=${period}`),
};

// Regulator APIs
export const regulatorAPI = {
  getDashboard: (period = '30d') => api.get(`/api/regulator/dashboard?period=${period}`),
  
  getStatistics: (period = '30d', clinicId = null) => {
    const params = { period };
    if (clinicId) params.clinic_id = clinicId;
    return api.get('/api/regulator/statistics', { params });
  },
  
  getClinics: (city, district) =>
    api.get('/api/regulator/clinics', { params: { city, district } }),
  
  getClinicDetails: (id, period = '30d') =>
    api.get(`/api/regulator/clinics/${id}?period=${period}`),
  
  getComplaints: (status) =>
    api.get('/api/regulator/complaints', { params: { status } }),
  
  getDiseaseAnalytics: (period = '30d') =>
    api.get(`/api/regulator/disease-analytics?period=${period}`),
};

export default api;
