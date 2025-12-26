import React, { useState, useEffect } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  Box,
  Alert,
} from '@mui/material';
import {
  MedicalServices,
  Assignment,
  LocalOffer,
  CalendarToday,
} from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { patientAPI } from '../../services/api';
import { useNavigate } from 'react-router-dom';

const PatientDashboard = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    scans: 0,
    plans: 0,
    appointments: 0,
  });
  const navigate = useNavigate();

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const [scansRes, plansRes, appointmentsRes] = await Promise.all([
        patientAPI.getScans(),
        patientAPI.getTreatmentPlans(),
        patientAPI.getAppointments(),
      ]);

      setStats({
        scans: scansRes.data.length,
        plans: plansRes.data.length,
        appointments: appointmentsRes.data.length,
      });
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Панель пациента">
        <LoadingSpinner />
      </Layout>
    );
  }

  const cards = [
    {
      title: 'Мои снимки КТ',
      value: stats.scans,
      icon: <MedicalServices fontSize="large" />,
      color: '#1976d2',
      action: () => navigate('/patient/scans'),
      buttonText: 'Просмотреть снимки',
    },
    {
      title: 'Планы лечения',
      value: stats.plans,
      icon: <Assignment fontSize="large" />,
      color: '#2e7d32',
      action: () => navigate('/patient/plans'),
      buttonText: 'Мои планы',
    },
    {
      title: 'Записи на приём',
      value: stats.appointments,
      icon: <CalendarToday fontSize="large" />,
      color: '#ed6c02',
      action: () => navigate('/patient/appointments'),
      buttonText: 'Мои записи',
    },
  ];

  return (
    <Layout title="Панель пациента">
      <Typography variant="h4" gutterBottom>
        Добро пожаловать!
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Управление вашими снимками, планами лечения и записями
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        У вас есть {stats.plans} активных плана лечения. Просмотрите предложения от клиник!
      </Alert>

      <Grid container spacing={3}>
        {cards.map((card, index) => (
          <Grid item xs={12} md={4} key={index}>
            <Card
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                transition: 'transform 0.2s',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: 4,
                },
              }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    mb: 2,
                  }}
                >
                  <Box
                    sx={{
                      backgroundColor: card.color,
                      color: 'white',
                      p: 1.5,
                      borderRadius: 2,
                    }}
                  >
                    {card.icon}
                  </Box>
                  <Typography variant="h3" component="div">
                    {card.value}
                  </Typography>
                </Box>
                <Typography variant="h6" gutterBottom>
                  {card.title}
                </Typography>
                <Button
                  variant="contained"
                  fullWidth
                  onClick={card.action}
                  sx={{ mt: 2 }}
                >
                  {card.buttonText}
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Layout>
  );
};

export default PatientDashboard;
