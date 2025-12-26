import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Box,
  Chip,
  Alert,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import {
  CheckCircle,
  Star,
  Payment,
  AccessTime,
} from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { patientAPI } from '../../services/api';
import { formatPrice } from '../../utils/constants';

const PatientOffers = () => {
  const { planId } = useParams();
  const [offers, setOffers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedOffer, setSelectedOffer] = useState(null);
  const [confirmDialog, setConfirmDialog] = useState(false);

  useEffect(() => {
    loadOffers();
  }, [planId]);

  const loadOffers = async () => {
    try {
      const response = await patientAPI.getOffers(planId);
      setOffers(response.data);
    } catch (error) {
      console.error('Failed to load offers:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectOffer = (offer) => {
    setSelectedOffer(offer);
    setConfirmDialog(true);
  };

  const handleConfirmSelection = async () => {
    try {
      await patientAPI.selectOffer(selectedOffer.id);
      setConfirmDialog(false);
      alert('Предложение принято! Запись создана.');
      loadOffers();
    } catch (error) {
      console.error('Failed to accept offer:', error);
      alert('Ошибка при принятии предложения');
    }
  };

  if (loading) {
    return (
      <Layout title="Предложения от клиник">
        <LoadingSpinner />
      </Layout>
    );
  }

  if (offers.length === 0) {
    return (
      <Layout title="Предложения от клиник">
        <Alert severity="info">
          Пока нет предложений от клиник. Подождите, клиники скоро направят свои предложения.
        </Alert>
      </Layout>
    );
  }

  return (
    <Layout title="Предложения от клиник">
      <Typography variant="h4" gutterBottom>
        Сравнение предложений
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Выберите наиболее подходящее предложение от клиник
      </Typography>

      <Grid container spacing={3}>
        {offers.map((offer) => (
          <Grid item xs={12} md={6} key={offer.id}>
            <Card
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                border: offer.status === 'accepted' ? '2px solid #4caf50' : 'none',
              }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                {/* Clinic Header */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                  <Box>
                    <Typography variant="h5" gutterBottom>
                      {offer.Clinic?.name}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Star sx={{ color: '#ffc107', fontSize: 20 }} />
                      <Typography variant="body2" color="text.secondary">
                        {offer.Clinic?.rating} ({offer.Clinic?.review_count} отзывов)
                      </Typography>
                    </Box>
                  </Box>
                  {offer.status === 'accepted' && (
                    <Chip
                      icon={<CheckCircle />}
                      label="Принято"
                      color="success"
                    />
                  )}
                </Box>

                {/* Clinic Info */}
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  📍 {offer.Clinic?.address}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  📞 {offer.Clinic?.phone || 'Не указано'}
                </Typography>

                {/* Price Breakdown */}
                <Box sx={{ my: 2 }}>
                  <Typography variant="h6" gutterBottom>
                    Стоимость по категориям:
                  </Typography>
                  <TableContainer>
                    <Table size="small">
                      <TableBody>
                        {offer.therapy_cost > 0 && (
                          <TableRow>
                            <TableCell>Терапия</TableCell>
                            <TableCell align="right">
                              {formatPrice(offer.therapy_cost)}
                            </TableCell>
                          </TableRow>
                        )}
                        {offer.orthopedics_cost > 0 && (
                          <TableRow>
                            <TableCell>Ортопедия</TableCell>
                            <TableCell align="right">
                              {formatPrice(offer.orthopedics_cost)}
                            </TableCell>
                          </TableRow>
                        )}
                        {offer.surgery_cost > 0 && (
                          <TableRow>
                            <TableCell>Хирургия</TableCell>
                            <TableCell align="right">
                              {formatPrice(offer.surgery_cost)}
                            </TableCell>
                          </TableRow>
                        )}
                        {offer.hygiene_cost > 0 && (
                          <TableRow>
                            <TableCell>Гигиена</TableCell>
                            <TableCell align="right">
                              {formatPrice(offer.hygiene_cost)}
                            </TableCell>
                          </TableRow>
                        )}
                        {offer.periodontics_cost > 0 && (
                          <TableRow>
                            <TableCell>Пародонтология</TableCell>
                            <TableCell align="right">
                              {formatPrice(offer.periodontics_cost)}
                            </TableCell>
                          </TableRow>
                        )}
                        <TableRow>
                          <TableCell>
                            <strong>Итого:</strong>
                          </TableCell>
                          <TableCell align="right">
                            <strong>{formatPrice(offer.total_cost)}</strong>
                          </TableCell>
                        </TableRow>
                      </TableBody>
                    </Table>
                  </TableContainer>
                </Box>

                {/* Additional Info */}
                <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                  <Chip
                    icon={<AccessTime />}
                    label={`Срок: ${offer.estimated_duration}`}
                    size="small"
                  />
                  {offer.installment_months > 0 && (
                    <Chip
                      icon={<Payment />}
                      label={`Рассрочка: ${offer.installment_months} мес`}
                      size="small"
                      color="primary"
                    />
                  )}
                </Box>

                {/* Warranty */}
                <Typography variant="body2" gutterBottom>
                  <strong>Гарантии:</strong> {offer.warranty_details}
                </Typography>

                {/* Notes */}
                {offer.notes && (
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    <strong>Примечания:</strong> {offer.notes}
                  </Typography>
                )}

                {/* Accept Button */}
                {offer.status !== 'accepted' && (
                  <Button
                    variant="contained"
                    fullWidth
                    size="large"
                    onClick={() => handleSelectOffer(offer)}
                    sx={{ mt: 2 }}
                  >
                    Выбрать это предложение
                  </Button>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Confirmation Dialog */}
      <Dialog open={confirmDialog} onClose={() => setConfirmDialog(false)}>
        <DialogTitle>Подтвердите выбор</DialogTitle>
        <DialogContent>
          <Typography>
            Вы уверены, что хотите выбрать предложение от клиники{' '}
            <strong>{selectedOffer?.Clinic?.name}</strong>?
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
            Стоимость: {formatPrice(selectedOffer?.total_cost)}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            После подтверждения будет создана запись на приём.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDialog(false)}>Отмена</Button>
          <Button onClick={handleConfirmSelection} variant="contained" autoFocus>
            Подтвердить
          </Button>
        </DialogActions>
      </Dialog>
    </Layout>
  );
};

export default PatientOffers;
