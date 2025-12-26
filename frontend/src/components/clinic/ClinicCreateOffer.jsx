import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Grid,
  Box,
  Alert,
} from '@mui/material';
import { Save } from '@mui/icons-material';
import Layout from '../common/Layout';
import { clinicAPI } from '../../services/api';

const ClinicCreateOffer = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const [formData, setFormData] = useState({
    therapy_cost: 0,
    orthopedics_cost: 0,
    surgery_cost: 0,
    hygiene_cost: 0,
    periodontics_cost: 0,
    estimated_duration: '',
    installment_months: 0,
    warranty_details: '',
    notes: '',
  });

  const handleChange = (field) => (event) => {
    const value = event.target.value;
    setFormData({ ...formData, [field]: value });
  };

  const calculateTotal = () => {
    return (
      parseInt(formData.therapy_cost || 0) +
      parseInt(formData.orthopedics_cost || 0) +
      parseInt(formData.surgery_cost || 0) +
      parseInt(formData.hygiene_cost || 0) +
      parseInt(formData.periodontics_cost || 0)
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const offerData = {
        treatment_plan_id: parseInt(planId),
        ...formData,
        therapy_cost: parseInt(formData.therapy_cost || 0),
        orthopedics_cost: parseInt(formData.orthopedics_cost || 0),
        surgery_cost: parseInt(formData.surgery_cost || 0),
        hygiene_cost: parseInt(formData.hygiene_cost || 0),
        periodontics_cost: parseInt(formData.periodontics_cost || 0),
        installment_months: parseInt(formData.installment_months || 0),
        total_cost: calculateTotal(),
      };

      await clinicAPI.createOffer(offerData);
      alert('Предложение успешно отправлено!');
      navigate('/clinic/incoming-plans');
    } catch (err) {
      console.error('Failed to create offer:', err);
      setError('Ошибка при создании предложения');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout title="Создать предложение">
      <Typography variant="h4" gutterBottom>
        Создать предложение
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Укажите стоимость и условия для плана лечения #{planId}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Стоимость по специализациям
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Терапия (₽)"
                  type="number"
                  value={formData.therapy_cost}
                  onChange={handleChange('therapy_cost')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Ортопедия (₽)"
                  type="number"
                  value={formData.orthopedics_cost}
                  onChange={handleChange('orthopedics_cost')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Хирургия (₽)"
                  type="number"
                  value={formData.surgery_cost}
                  onChange={handleChange('surgery_cost')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Гигиена (₽)"
                  type="number"
                  value={formData.hygiene_cost}
                  onChange={handleChange('hygiene_cost')}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Пародонтология (₽)"
                  type="number"
                  value={formData.periodontics_cost}
                  onChange={handleChange('periodontics_cost')}
                />
              </Grid>
            </Grid>

            <Box sx={{ mt: 3, p: 2, bgcolor: 'primary.main', color: 'white', borderRadius: 1 }}>
              <Typography variant="h6">
                Общая стоимость: {calculateTotal().toLocaleString('ru-RU')} ₽
              </Typography>
            </Box>
          </CardContent>
        </Card>

        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Дополнительная информация
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Срок лечения"
                  placeholder="например: 2-3 месяца"
                  value={formData.estimated_duration}
                  onChange={handleChange('estimated_duration')}
                  required
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Рассрочка (месяцев)"
                  type="number"
                  value={formData.installment_months}
                  onChange={handleChange('installment_months')}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  multiline
                  rows={2}
                  label="Гарантии"
                  placeholder="Укажите условия гарантии"
                  value={formData.warranty_details}
                  onChange={handleChange('warranty_details')}
                  required
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  multiline
                  rows={3}
                  label="Примечания"
                  placeholder="Дополнительная информация"
                  value={formData.notes}
                  onChange={handleChange('notes')}
                />
              </Grid>
            </Grid>
          </CardContent>
        </Card>

        <Button
          type="submit"
          variant="contained"
          size="large"
          fullWidth
          startIcon={<Save />}
          disabled={loading}
        >
          {loading ? 'Отправка...' : 'Отправить предложение'}
        </Button>
      </form>
    </Layout>
  );
};

export default ClinicCreateOffer;
