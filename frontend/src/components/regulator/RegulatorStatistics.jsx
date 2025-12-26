import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Grid,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
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

const RegulatorStatistics = () => {
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('30d');
  const [statistics, setStatistics] = useState([]);

  useEffect(() => {
    loadStatistics();
  }, [period]);

  const loadStatistics = async () => {
    try {
      const response = await regulatorAPI.getStatistics(period);
      setStatistics(response.data.statistics || []);
    } catch (error) {
      console.error('Failed to load statistics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Статистика">
        <LoadingSpinner />
      </Layout>
    );
  }

  // Prepare chart data
  const chartData = statistics.slice(-30).map((stat) => ({
    date: new Date(stat.date).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' }),
    plans: stat.treatment_plans_generated,
    appointments: stat.appointments_completed,
    revenue: stat.total_revenue / 1000, // В тысячах для читаемости
    patients: stat.patient_count,
  }));

  // Calculate totals
  const totals = statistics.reduce(
    (acc, stat) => ({
      plans: acc.plans + stat.treatment_plans_generated,
      appointments: acc.appointments + stat.appointments_completed,
      revenue: acc.revenue + stat.total_revenue,
      patients: acc.patients + stat.patient_count,
    }),
    { plans: 0, appointments: 0, revenue: 0, patients: 0 }
  );

  return (
    <Layout title="Статистика">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4">
          Детальная статистика
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

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h4" color="primary">
                {totals.plans}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Всего планов
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h4" color="success.main">
                {totals.appointments}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Завершённых приёмов
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h4" color="warning.main">
                {totals.patients}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Пациентов
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h4" color="secondary.main">
                {formatPrice(totals.revenue)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Общая выручка
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Динамика планов и приёмов
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={chartData}>
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

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Выручка (тыс. ₽)
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Area type="monotone" dataKey="revenue" stroke="#9c27b0" fill="#9c27b0" fillOpacity={0.3} name="Выручка" />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Layout>
  );
};

export default RegulatorStatistics;
