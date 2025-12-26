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
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material';
import { LocalOffer } from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { clinicAPI } from '../../services/api';
import { formatPrice, SPECIALIZATIONS } from '../../utils/constants';

const ClinicIncomingPlans = () => {
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadPlans();
  }, []);

  const loadPlans = async () => {
    try {
      const response = await clinicAPI.getIncomingPlans();
      setPlans(response.data);
    } catch (error) {
      console.error('Failed to load plans:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Входящие планы лечения">
        <LoadingSpinner />
      </Layout>
    );
  }

  return (
    <Layout title="Входящие планы лечения">
      <Typography variant="h4" gutterBottom>
        Входящие планы лечения
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Новые планы лечения от пациентов, ожидающие вашего предложения
      </Typography>

      {plans.length === 0 ? (
        <Card>
          <CardContent>
            <Typography variant="body1" color="text.secondary" align="center">
              Нет новых планов лечения
            </Typography>
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={3}>
          {plans.map((plan) => (
            <Grid item xs={12} key={plan.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                    <Typography variant="h6">
                      План лечения #{plan.id}
                    </Typography>
                    <Chip label={plan.status === 'generated' ? 'Новый' : plan.status} color="primary" />
                  </Box>

                  <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
                    {plan.requires_therapy && <Chip label="Терапия" size="small" />}
                    {plan.requires_orthopedics && <Chip label="Ортопедия" size="small" />}
                    {plan.requires_surgery && <Chip label="Хирургия" size="small" />}
                    {plan.requires_hygiene && <Chip label="Гигиена" size="small" />}
                    {plan.requires_periodontics && <Chip label="Пародонтология" size="small" />}
                  </Box>

                  <TableContainer sx={{ mb: 2 }}>
                    <Table size="small">
                      <TableHead>
                        <TableRow>
                          <TableCell>Специализация</TableCell>
                          <TableCell align="right">Мин. стоимость</TableCell>
                          <TableCell align="right">Макс. стоимость</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {plan.therapy_min_cost > 0 && (
                          <TableRow>
                            <TableCell>Терапия</TableCell>
                            <TableCell align="right">{formatPrice(plan.therapy_min_cost)}</TableCell>
                            <TableCell align="right">{formatPrice(plan.therapy_max_cost)}</TableCell>
                          </TableRow>
                        )}
                        {plan.orthopedics_min_cost > 0 && (
                          <TableRow>
                            <TableCell>Ортопедия</TableCell>
                            <TableCell align="right">{formatPrice(plan.orthopedics_min_cost)}</TableCell>
                            <TableCell align="right">{formatPrice(plan.orthopedics_max_cost)}</TableCell>
                          </TableRow>
                        )}
                        {plan.surgery_min_cost > 0 && (
                          <TableRow>
                            <TableCell>Хирургия</TableCell>
                            <TableCell align="right">{formatPrice(plan.surgery_min_cost)}</TableCell>
                            <TableCell align="right">{formatPrice(plan.surgery_max_cost)}</TableCell>
                          </TableRow>
                        )}
                        {plan.hygiene_min_cost > 0 && (
                          <TableRow>
                            <TableCell>Гигиена</TableCell>
                            <TableCell align="right">{formatPrice(plan.hygiene_min_cost)}</TableCell>
                            <TableCell align="right">{formatPrice(plan.hygiene_max_cost)}</TableCell>
                          </TableRow>
                        )}
                      </TableBody>
                    </Table>
                  </TableContainer>

                  <Button
                    variant="contained"
                    fullWidth
                    startIcon={<LocalOffer />}
                    onClick={() => navigate(`/clinic/create-offer/${plan.id}`)}
                  >
                    Создать предложение
                  </Button>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Layout>
  );
};

export default ClinicIncomingPlans;
