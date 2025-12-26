import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Chip,
  Box,
  TextField,
  InputAdornment,
} from '@mui/material';
import {
  Search,
  Star,
  Visibility,
  CheckCircle,
} from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { regulatorAPI } from '../../services/api';

const RegulatorClinics = () => {
  const [clinics, setClinics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    loadClinics();
  }, []);

  const loadClinics = async () => {
    try {
      const response = await regulatorAPI.getClinics();
      setClinics(response.data);
    } catch (error) {
      console.error('Failed to load clinics:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredClinics = clinics.filter((clinic) =>
    clinic.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    clinic.city?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    clinic.district?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (loading) {
    return (
      <Layout title="Клиники">
        <LoadingSpinner />
      </Layout>
    );
  }

  return (
    <Layout title="Клиники">
      <Typography variant="h4" gutterBottom>
        Все клиники региона
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
        Полный список зарегистрированных стоматологических клиник
      </Typography>

      <TextField
        fullWidth
        placeholder="Поиск по названию, городу или району..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        sx={{ mb: 3 }}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <Search />
            </InputAdornment>
          ),
        }}
      />

      <Grid container spacing={3}>
        {filteredClinics.map((clinic) => (
          <Grid item xs={12} md={6} key={clinic.id}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                  <Typography variant="h6">
                    {clinic.name}
                  </Typography>
                  {clinic.is_active && (
                    <Chip
                      icon={<CheckCircle />}
                      label="Активна"
                      color="success"
                      size="small"
                    />
                  )}
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                  <Star sx={{ color: '#ffc107', fontSize: 20 }} />
                  <Typography variant="body2" color="text.secondary">
                    {clinic.rating?.toFixed(1)} ({clinic.review_count} отзывов)
                  </Typography>
                </Box>

                <Typography variant="body2" color="text.secondary" gutterBottom>
                  📍 {clinic.address}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  🏙️ {clinic.city}, {clinic.district}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  📋 Лицензия: {clinic.license_number}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  📅 Основана: {clinic.year_established}
                </Typography>

                <Box sx={{ display: 'flex', gap: 1, mt: 2, flexWrap: 'wrap' }}>
                  {clinic.has_therapy && <Chip label="Терапия" size="small" />}
                  {clinic.has_orthopedics && <Chip label="Ортопедия" size="small" />}
                  {clinic.has_surgery && <Chip label="Хирургия" size="small" />}
                  {clinic.has_hygiene && <Chip label="Гигиена" size="small" />}
                  {clinic.has_periodontics && <Chip label="Пародонтология" size="small" />}
                </Box>

                <Box sx={{ display: 'flex', gap: 1, mt: 2 }}>
                  {clinic.offers_installment && (
                    <Chip label="Рассрочка" color="primary" size="small" />
                  )}
                  {clinic.offers_insurance && (
                    <Chip label="Страховка" color="primary" size="small" />
                  )}
                </Box>

                <Button
                  variant="outlined"
                  fullWidth
                  startIcon={<Visibility />}
                  onClick={() => navigate(`/regulator/clinics/${clinic.id}`)}
                  sx={{ mt: 2 }}
                >
                  Подробная информация
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {filteredClinics.length === 0 && (
        <Card>
          <CardContent>
            <Typography variant="body1" color="text.secondary" align="center">
              Клиники не найдены
            </Typography>
          </CardContent>
        </Card>
      )}
    </Layout>
  );
};

export default RegulatorClinics;
