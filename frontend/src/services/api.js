// Get API URL from environment variable or use default
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Common API functions
export const commonAPI = {
  getConstants: async () => {
    const response = await fetch(`${API_URL}/api/constants`);
    if (!response.ok) {
      throw new Error('Failed to fetch constants');
    }
    return response.json();
  }
};

// Auth API functions
export const authAPI = {
  login: async (username, password) => {
    const response = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Login failed');
    }

    return response.json();
  },

  refreshToken: async (refreshToken) => {
    const response = await fetch(`${API_URL}/api/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      throw new Error('Token refresh failed');
    }

    return response.json();
  }
};

// Helper function to make authenticated requests
const authenticatedFetch = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    if (response.status === 401) {
      // Token expired, redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }
    const error = await response.json();
    throw new Error(error.message || 'Request failed');
  }

  return response.json();
};

// Patient API functions
export const patientAPI = {
  getScans: () => authenticatedFetch('/api/patient/scans'),
  getTreatmentPlans: () => authenticatedFetch('/api/patient/plans'),
  getOffers: (planId) => authenticatedFetch(`/api/patient/plans/${planId}/offers`),
  
  // New methods for appointments
  getAppointments: () => authenticatedFetch('/api/patient/appointments'),
  bookAppointment: (appointmentData) => authenticatedFetch('/api/patient/appointments', {
    method: 'POST',
    body: JSON.stringify(appointmentData),
  }),
  getClinics: () => authenticatedFetch('/api/patient/clinics'),
  getProfile: () => authenticatedFetch('/api/patient/profile'),
};

// Clinic API functions
export const clinicAPI = {
  getDashboard: () => authenticatedFetch('/api/clinic/dashboard'),
  getIncomingPlans: () => authenticatedFetch('/api/clinic/incoming-plans'),
  getPriceList: () => authenticatedFetch('/api/clinic/price-list'),
  
  // New methods for appointments
  getAppointments: () => authenticatedFetch('/api/clinic/appointments'),
  updateAppointmentStatus: (appointmentId, status) => authenticatedFetch(`/api/clinic/appointments/${appointmentId}`, {
    method: 'PATCH',
    body: JSON.stringify({ status }),
  }),
  getProfile: () => authenticatedFetch('/api/clinic/profile'),
  updateProfile: (profileData) => authenticatedFetch('/api/clinic/profile', {
    method: 'PUT',
    body: JSON.stringify(profileData),
  }),
};

// Regulator API functions
export const regulatorAPI = {
  getDashboard: () => authenticatedFetch('/api/regulator/dashboard'),
  getStatistics: () => authenticatedFetch('/api/regulator/statistics'),
  getClinics: () => authenticatedFetch('/api/regulator/clinics'),
  
  // New methods for appointments
  getAppointments: () => authenticatedFetch('/api/regulator/appointments'),
  updateClinicStatus: (clinicId, status) => authenticatedFetch(`/api/regulator/clinics/${clinicId}`, {
    method: 'PATCH',
    body: JSON.stringify({ status }),
  }),
};
