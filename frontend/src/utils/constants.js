// Utility functions for formatting and display

export const formatPrice = (price) => {
  if (price === null || price === undefined) return '—';
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price);
};

export const formatDate = (dateString) => {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);
};

export const formatDateTime = (dateString) => {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

export const getStatusColor = (status) => {
  const statusColors = {
    generated: '#2196F3',
    offers_received: '#FF9800',
    offer_selected: '#4CAF50',
    in_progress: '#9C27B0',
    completed: '#4CAF50',
    cancelled: '#F44336',
    pending: '#FFC107',
    processing: '#2196F3',
    failed: '#F44336',
    sent: '#2196F3',
    viewed: '#FF9800',
    accepted: '#4CAF50',
    rejected: '#F44336',
  };
  return statusColors[status] || '#757575';
};

export const getUrgencyColor = (urgency) => {
  const urgencyColors = {
    high: '#F44336',
    medium: '#FF9800',
    low: '#4CAF50',
  };
  return urgencyColors[urgency] || '#757575';
};

// Export empty objects for backward compatibility
export const ROLES = {};
export const SPECIALIZATIONS = {};
export const TREATMENT_STATUSES = {};
export const SCAN_STATUSES = {};
export const OFFER_STATUSES = {};
export const CITIES = [];
export const DISTRICTS_BY_CITY = {};
export const PRICE_SEGMENTS = [];
