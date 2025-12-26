import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ConstantsProvider } from './contexts/ConstantsContext';

import Login from './pages/Login';
import PatientDashboard from './components/patient/PatientDashboard';
import ClinicDashboard from './components/clinic/ClinicDashboard';
import RegulatorDashboard from './components/regulator/RegulatorDashboard';

function App() {
  return (
    <AuthProvider>
      <ConstantsProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/patient/dashboard" element={<PatientDashboard />} />
            <Route path="/clinic/dashboard" element={<ClinicDashboard />} />
            <Route path="/regulator/dashboard" element={<RegulatorDashboard />} />
            <Route path="/" element={<Navigate to="/login" replace />} />
          </Routes>
        </Router>
      </ConstantsProvider>
    </AuthProvider>
  );
}

export default App;
