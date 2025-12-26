import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Grid,
} from '@mui/material';
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  LineChart,
  Line,
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

const COLORS = ['#1976d2', '#2e7d32', '#ed6c02', '#9c27b0', '#d32f2f'];

const RegulatorDiseaseAnalytics = () => {
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('30d');
  const [analytics, setAnalytics] = useState(null);

  useEffect(() => {
    loadAnalytics();
  }, [period]);

  const loadAnalytics = async () => {
    try {
      const response = await regulatorAPI.getDiseaseAnalytics(period);
      setAnalytics(response.data);
    } catch (error) {
      console.error('Failed to load analytics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Аналитика заболеваний">
        <LoadingSpinner />
      </Layout>
    );
  }

  const diseases = analytics?.diseases || [];
  const timeSeries = analytics?.time_series?.slice(-30) || [];

  // Prepare time series data
  const timeSeriesData = timeSeries.map((item) => ({
    date: new Date(item.date).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' }),
    caries: item.caries_count,
    pulpitis: item.pulpitis_count,
    periodontitis: item.periodontitis_count,
    gingivitis: item.gingivitis_count,
    parodontitis: item.parodontitis_count,
  }));

  return (
    <Layout title="Аналитика заболеваний">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4">
          Аналитика заболеваний
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

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Общая статистика
          </Typography>
          <Typography variant="h4" color="primary">
            {analytics?.total_cases || 0} случаев
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Всего зарегистрировано случаев за период
          </Typography>
        </CardContent>
      </Card>

      <Grid container spacing={3}>
        {/* Pie Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Распределение заболеваний
              </Typography>
              <ResponsiveContainer width="100%" height={400}>
                <PieChart>
                  <Pie
                    data={diseases}
                    dataKey="count"
                    nameKey="disease"
                    cx="50%"
                    cy="50%"
                    outerRadius={120}
                    label={(entry) => `${entry.disease}: ${entry.count}`}
                  >
                    {diseases.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Bar Chart */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Количество случаев
              </Typography>
              <ResponsiveContainer width="100%" height={400}>
                <BarChart data={diseases}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="disease" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="count" fill="#1976d2" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Time Series Chart */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Динамика заболеваний за период
              </Typography>
              <ResponsiveContainer width="100%" height={400}>
                <LineChart data={timeSeriesData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Line type="monotone" dataKey="caries" stroke="#1976d2" name="Кариес" />
                  <Line type="monotone" dataKey="pulpitis" stroke="#2e7d32" name="Пульпит" />
                  <Line type="monotone" dataKey="periodontitis" stroke="#ed6c02" name="Периодонтит" />
                  <Line type="monotone" dataKey="gingivitis" stroke="#9c27b0" name="Гингивит" />
                  <Line type="monotone" dataKey="parodontitis" stroke="#d32f2f" name="Пародонтит" />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Statistics Table */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Детальная статистика
              </Typography>
              <Grid container spacing={2}>
                {diseases.map((disease, index) => (
                  <Grid item xs={12} sm={6} md={4} key={index}>
                    <Box
                      sx={{
                        p: 2,
                        border: 1,
                        borderColor: 'divider',
                        borderRadius: 1,
                        backgroundColor: COLORS[index % COLORS.length] + '10',
                      }}
                    >
                      <Typography variant="h5" sx={{ color: COLORS[index % COLORS.length] }}>
                        {disease.count}
                      </Typography>
                      <Typography variant="body1" fontWeight="bold">
                        {disease.disease}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {disease.percentage?.toFixed(1)}% от общего числа
                      </Typography>
                    </Box>
                  </Grid>
                ))}
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Layout>
  );
};

export default RegulatorDiseaseAnalytics;
