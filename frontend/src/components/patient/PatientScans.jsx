import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Chip,
  Box,
} from '@mui/material';
import {
  CheckCircle,
  Schedule,
  Visibility,
} from '@mui/icons-material';
import Layout from '../common/Layout';
import LoadingSpinner from '../common/LoadingSpinner';
import { patientAPI } from '../../services/api';
import { useNavigate } from 'react-router-dom';
import { formatDateTime } from '../../utils/constants';

const PatientScans = () => {
  const [scans, setScans] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadScans();
  }, []);

  const loadScans = async () => {
    try {
      const response = await patientAPI.getScans();
      setScans(response.data);
    } catch (error) {
      console.error('Failed to load scans:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout title="Мои снимки КТ">
        <LoadingSpinner />
      </Layout>
    );
  }

  return (
    <Layout title="Мои снимки КТ">
      <Typography variant="h4" gutterBottom>
        Мои снимки КТ
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Все загруженные снимки и результаты анализа ИИ
      </Typography>

      {scans.length === 0 ? (
        <Card>
          <CardContent>
            <Typography variant="body1" color="text.secondary" align="center">
              У вас пока нет загруженных снимков
            </Typography>
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={3}>
          {scans.map((scan) => (
            <Grid item xs={12} md={6} key={scan.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                    <Typography variant="h6">
                      Снимок #{scan.id}
                    </Typography>
                    {scan.ai_processed ? (
                      <Chip
                        icon={<CheckCircle />}
                        label="Обработано"
                        color="success"
                        size="small"
                      />
                    ) : (
                      <Chip
                        icon={<Schedule />}
                        label="В обработке"
                        color="warning"
                        size="small"
                      />
                    )}
                  </Box>

                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Дата загрузки: {formatDateTime(scan.upload_date)}
                  </Typography>

                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Статус: {scan.status === 'completed' ? 'Завершено' : 'В процессе'}
                  </Typography>

                  {scan.ai_processed && (
                    <Button
                      variant="contained"
                      fullWidth
                      startIcon={<Visibility />}
                      onClick={() => navigate(`/patient/scans/${scan.id}`)}
                      sx={{ mt: 2 }}
                    >
                      Просмотреть план лечения
                    </Button>
                  )}
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Layout>
  );
};

export default PatientScans;
