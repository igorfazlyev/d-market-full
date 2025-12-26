import React, { useState, useEffect } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import {
  TrendingUp,
  People,
  Assignment,
  AttachMoney,
} from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { clinicAPI } from '../../services/api';
import { formatPrice } from '../../utils/constants';

const ClinicDashboard = () => {
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('30d');
  const [metrics, setMetrics] = useState(null);

  useEffect(() => {
    loadDashboard();
  }, [period]);

  const loadDashboard = async () => {
    try {
      const response = await clinicAPI.getDashboard(period);
      setMetrics(response.data);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Панель клиники">
        <LoadingSpinner />
      </Layout>
    );
  }

  const cards = [
    {
      title: 'Всего пациентов',
      value: metrics?.patient_count || 0,
      icon: <People fontSize="large" />,
      color: '#1976d2',
    },
    {
      title: 'Планов лечения',
      value: metrics?.treatment_plans_generated || 0,
      icon: <Assignment fontSize="large" />,
      color: '#2e7d32',
    },
    {
      title: 'Завершённых приёмов',
      value: metrics?.appointments_completed || 0,
      icon: <TrendingUp fontSize="large" />,
      color: '#ed6c02',
    },
    {
      title: 'Выручка',
      value: formatPrice(metrics?.total_revenue || 0),
      icon: <AttachMoney fontSize="large" />,
      color: '#9c27b0',
    },
  ];

  return (
    <Layout title="Панель клиники">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4">
          Панель управления клиники
        </Typography>
        <FormControl sx={{ minWidth: 150 }}>
          <InputLabel>Период</InputLabel>
          <Select
            value={period}
            label="Период"
            onChange={(e) => setPeriod(e.target.value)}
          >
            <MenuItem value="7d">7 дней</MenuItem>
            <MenuItem value="30d">30 дней</MenuItem>
            <MenuItem value="90d">90 дней</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Grid container spacing={3}>
        {cards.map((card, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Card>
              <CardContent>
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
                </Box>
                <Typography variant="h4" component="div" gutterBottom>
                  {card.value}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {card.title}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3} sx={{ mt: 2 }}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Средняя стоимость лечения
              </Typography>
              <Typography variant="h4">
                {formatPrice(metrics?.average_treatment_cost || 0)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Среднее время ожидания
              </Typography>
              <Typography variant="h4">
                {metrics?.average_wait_days?.toFixed(1) || 0} дней
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Layout>
  );
};

export default ClinicDashboard;
