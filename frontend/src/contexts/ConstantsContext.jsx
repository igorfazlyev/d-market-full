import React, { createContext, useContext, useState, useEffect } from 'react';
import { commonAPI } from '../services/api';

const ConstantsContext = createContext();

export const useConstants = () => {
  const context = useContext(ConstantsContext);
  if (!context) {
    throw new Error('useConstants must be used within ConstantsProvider');
  }
  return context;
};

export const ConstantsProvider = ({ children }) => {
  const [constants, setConstants] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchConstants = async () => {
      try {
        const data = await commonAPI.getConstants();
        setConstants(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch constants:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchConstants();
  }, []);

  if (loading) {
    return <div>Loading constants...</div>;
  }

  if (error) {
    return <div>Error loading constants: {error}</div>;
  }

  return (
    <ConstantsContext.Provider value={constants}>
      {children}
    </ConstantsContext.Provider>
  );
};
