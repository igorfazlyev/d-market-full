import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './components/common/PrivateRoute';
import Login from './pages/Login';
import Unauthorized from './pages/Unauthorized';

// Patient Components
import PatientDashboard from './components/patient/PatientDashboard';
import PatientScans from './components/patient/PatientScans';
import PatientTreatmentPlan from './components/patient/PatientTreatmentPlan';
import PatientOffers from './components/patient/PatientOffers';

// Clinic Components
import ClinicDashboard from './components/clinic/ClinicDashboard';
import ClinicIncomingPlans from './components/clinic/ClinicIncomingPlans';
import ClinicCreateOffer from './components/clinic/ClinicCreateOffer';

// Regulator Components
import RegulatorDashboard from './components/regulator/RegulatorDashboard';
import RegulatorClinics from './components/regulator/RegulatorClinics';
import RegulatorStatistics from './components/regulator/RegulatorStatistics';
import RegulatorDiseaseAnalytics from './components/regulator/RegulatorDiseaseAnalytics';

import { ROLES } from './utils/constants';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
    },
    h5: {
      fontWeight: 600,
    },
    h6: {
      fontWeight: 600,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Routes>
            {/* Public Routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/unauthorized" element={<Unauthorized />} />

            {/* Patient Routes */}
            <Route
              path="/patient/dashboard"
              element={
                <PrivateRoute allowedRoles={[ROLES.PATIENT]}>
                  <PatientDashboard />
                </PrivateRoute>
              }
            />
            <Route
              path="/patient/scans"
              element={
                <PrivateRoute allowedRoles={[ROLES.PATIENT]}>
                  <PatientScans />
                </PrivateRoute>
              }
            />
            <Route
              path="/patient/scans/:id"
              element={
                <PrivateRoute allowedRoles={[ROLES.PATIENT]}>
                  <PatientTreatmentPlan />
                </PrivateRoute>
              }
            />
            <Route
              path="/patient/plans/:planId/offers"
              element={
                <PrivateRoute allowedRoles={[ROLES.PATIENT]}>
                  <PatientOffers />
                </PrivateRoute>
              }
            />

            {/* Clinic Routes */}
            <Route
              path="/clinic/dashboard"
              element={
                <PrivateRoute allowedRoles={[ROLES.CLINIC]}>
                  <ClinicDashboard />
                </PrivateRoute>
              }
            />
            <Route
              path="/clinic/incoming-plans"
              element={
                <PrivateRoute allowedRoles={[ROLES.CLINIC]}>
                  <ClinicIncomingPlans />
                </PrivateRoute>
              }
            />
            <Route
              path="/clinic/create-offer/:planId"
              element={
                <PrivateRoute allowedRoles={[ROLES.CLINIC]}>
                  <ClinicCreateOffer />
                </PrivateRoute>
              }
            />

            {/* Regulator Routes */}
            <Route
              path="/regulator/dashboard"
              element={
                <PrivateRoute allowedRoles={[ROLES.REGULATOR]}>
                  <RegulatorDashboard />
                </PrivateRoute>
              }
            />
            <Route
              path="/regulator/clinics"
              element={
                <PrivateRoute allowedRoles={[ROLES.REGULATOR]}>
                  <RegulatorClinics />
                </PrivateRoute>
              }
            />
            <Route
              path="/regulator/statistics"
              element={
                <PrivateRoute allowedRoles={[ROLES.REGULATOR]}>
                  <RegulatorStatistics />
                </PrivateRoute>
              }
            />
            <Route
              path="/regulator/disease-analytics"
              element={
                <PrivateRoute allowedRoles={[ROLES.REGULATOR]}>
                  <RegulatorDiseaseAnalytics />
                </PrivateRoute>
              }
            />

            {/* Default Route */}
            <Route path="/" element={<Navigate to="/login" replace />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
