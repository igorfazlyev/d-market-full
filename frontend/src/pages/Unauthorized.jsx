import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
} from '@mui/material';
import { Block } from '@mui/icons-material';

const Unauthorized = () => {
  const navigate = useNavigate();

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Paper elevation={3} sx={{ p: 4, width: '100%', textAlign: 'center' }}>
          <Block sx={{ fontSize: 80, color: 'error.main', mb: 2 }} />
          <Typography variant="h4" gutterBottom>
            Доступ запрещён
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            У вас нет прав для просмотра этой страницы
          </Typography>
          <Button
            variant="contained"
            size="large"
            onClick={() => navigate('/login')}
          >
            Вернуться к входу
          </Button>
        </Paper>
      </Box>
    </Container>
  );
};

export default Unauthorized;
