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
  Business,
  Assignment,
  People,
  AttachMoney,
  TrendingUp,
  AccessTime,
} from '@mui/icons-material';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { regulatorAPI } from '../../services/api';
import { formatPrice } from '../../utils/constants';

const RegulatorDashboard = () => {
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('30d');
  const [dashboardData, setDashboardData] = useState(null);

  useEffect(() => {
    loadDashboard();
  }, [period]);

  const loadDashboard = async () => {
    try {
      const response = await regulatorAPI.getDashboard(period);
      setDashboardData(response.data);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Панель регулятора">
        <LoadingSpinner />
      </Layout>
    );
  }

  const summary = dashboardData?.summary || {};
  const diseases = dashboardData?.disease_statistics || {};

  const statsCards = [
    {
      title: 'Всего клиник',
      value: summary.total_clinics || 0,
      icon: <Business fontSize="large" />,
      color: '#1976d2',
    },
    {
      title: 'Планов лечения',
      value: summary.total_treatment_plans || 0,
      icon: <Assignment fontSize="large" />,
      color: '#2e7d32',
    },
    {
      title: 'Всего пациентов',
      value: summary.total_patients || 0,
      icon: <People fontSize="large" />,
      color: '#ed6c02',
    },
    {
      title: 'Общая выручка',
      value: formatPrice(summary.total_revenue || 0),
      icon: <AttachMoney fontSize="large" />,
      color: '#9c27b0',
    },
    {
      title: 'Среднее время ожидания',
      value: `${summary.average_wait_days?.toFixed(1) || 0} дней`,
      icon: <AccessTime fontSize="large" />,
      color: '#f57c00',
    },
    {
      title: 'Средняя стоимость лечения',
      value: formatPrice(summary.average_treatment_cost || 0),
      icon: <TrendingUp fontSize="large" />,
      color: '#0288d1',
    },
  ];

  // Prepare disease data for chart
  const diseaseData = [
    { name: 'Кариес', value: diseases.caries || 0 },
    { name: 'Пульпит', value: diseases.pulpitis || 0 },
    { name: 'Периодонтит', value: diseases.periodontitis || 0 },
    { name: 'Гингивит', value: diseases.gingivitis || 0 },
    { name: 'Пародонтит', value: diseases.parodontitis || 0 },
  ];

  // Prepare time series data for chart
  const timeSeriesData = dashboardData?.time_series?.slice(-30).map((item) => ({
    date: new Date(item.date).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' }),
    plans: item.treatment_plans_generated,
    appointments: item.appointments_completed,
  })) || [];

  return (
    <Layout title="Панель регулятора">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4">
          Региональная панель мониторинга
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

      {/* Statistics Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statsCards.map((card, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
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
                <Typography variant="h5" component="div" gutterBottom>
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

      {/* Charts */}
      <Grid container spacing={3}>
        {/* Disease Statistics Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Статистика заболеваний
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={diseaseData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="value" fill="#1976d2" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Time Series Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Динамика за период
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={timeSeriesData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Line type="monotone" dataKey="plans" stroke="#1976d2" name="Планы лечения" />
                  <Line type="monotone" dataKey="appointments" stroke="#2e7d32" name="Приёмы" />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Layout>
  );
};

export default RegulatorDashboard;
