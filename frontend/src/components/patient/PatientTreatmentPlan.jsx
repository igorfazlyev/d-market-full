import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  Box,
  Alert,
} from '@mui/material';
import { LocalOffer } from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { patientAPI } from '../../services/api';
import { formatPrice, SPECIALIZATIONS } from '../../utils/constants';

const PatientTreatmentPlan = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [plan, setPlan] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPlan();
  }, [id]);

  const loadPlan = async () => {
    try {
      const response = await patientAPI.getTreatmentPlan(id);
      setPlan(response.data);
    } catch (error) {
      console.error('Failed to load plan:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="План лечения">
        <LoadingSpinner />
      </Layout>
    );
  }

  if (!plan) {
    return (
      <Layout title="План лечения">
        <Alert severity="error">План лечения не найден</Alert>
      </Layout>
    );
  }

  const getUrgencyColor = (urgency) => {
    switch (urgency) {
      case 'high':
        return 'error';
      case 'medium':
        return 'warning';
      case 'low':
        return 'success';
      default:
        return 'default';
    }
  };

  const getUrgencyText = (urgency) => {
    switch (urgency) {
      case 'high':
        return 'Высокая';
      case 'medium':
        return 'Средняя';
      case 'low':
        return 'Низкая';
      default:
        return urgency;
    }
  };

  return (
    <Layout title="План лечения">
      <Typography variant="h4" gutterBottom>
        План лечения
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
        Детальный план лечения на основе анализа КТ снимка
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        Этот план создан на основе анализа вашего КТ снимка с использованием ИИ.
        Окончательные рекомендации будут даны стоматологом при осмотре.
      </Alert>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Необходимые специализации
          </Typography>
          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
            {plan.requires_therapy && <Chip label="Терапия" color="primary" />}
            {plan.requires_orthopedics && <Chip label="Ортопедия" color="primary" />}
            {plan.requires_surgery && <Chip label="Хирургия" color="primary" />}
            {plan.requires_hygiene && <Chip label="Гигиена" color="primary" />}
            {plan.requires_periodontics && <Chip label="Пародонтология" color="primary" />}
          </Box>
        </CardContent>
      </Card>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Процедуры
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Зуб</TableCell>
                  <TableCell>Специализация</TableCell>
                  <TableCell>Диагноз</TableCell>
                  <TableCell>Процедура</TableCell>
                  <TableCell>Срочность</TableCell>
                  <TableCell align="right">Стоимость</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {plan.Items?.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>{item.tooth_number}</TableCell>
                    <TableCell>{SPECIALIZATIONS[item.specialization]}</TableCell>
                    <TableCell>{item.diagnosis}</TableCell>
                    <TableCell>{item.procedure}</TableCell>
                    <TableCell>
                      <Chip
                        label={getUrgencyText(item.urgency)}
                        color={getUrgencyColor(item.urgency)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right">
                      {formatPrice(item.estimated_cost)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      <Button
        variant="contained"
        size="large"
        fullWidth
        startIcon={<LocalOffer />}
        onClick={() => navigate(`/patient/plans/${plan.id}/offers`)}
      >
        Просмотреть предложения от клиник
      </Button>
    </Layout>
  );
};

export default PatientTreatmentPlan;
