import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  Paper,
  List,
  ListItem,
  ListItemText,
  Divider,
} from '@mui/material';
import { Login as LoginIcon } from '@mui/icons-material';
import { authAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await authAPI.login(username, password);
      login(response.access_token, response.user);
      
      switch (response.user.role) {
        case 'patient':
          navigate('/patient/dashboard');
          break;
        case 'clinic':
          navigate('/clinic/dashboard');
          break;
        case 'regulator':
          navigate('/regulator/dashboard');
          break;
        default:
          navigate('/');
      }
    } catch (err) {
      setError(err.message || 'Ошибка входа');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#f5f7fa',
        padding: 2,
      }}
    >
      <Container maxWidth="sm">
        <Box sx={{ textAlign: 'center', mb: 4 }}>
          <Typography
            variant="h3"
            component="h1"
            sx={{
              color: '#2c3e50',
              fontWeight: 700,
              mb: 1,
            }}
          >
            🦷 Dental Marketplace
          </Typography>
          <Typography variant="h6" sx={{ color: '#7f8c8d' }}>
            Платформа для стоматологических услуг
          </Typography>
        </Box>

        <Card elevation={2} sx={{ borderRadius: 3, backgroundColor: 'white' }}>
          <CardContent sx={{ p: 4 }}>
            <Box component="form" onSubmit={handleSubmit}>
              {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {error}
                </Alert>
              )}

              <TextField
                fullWidth
                label="Имя пользователя"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                margin="normal"
                required
                autoFocus
                sx={{
                  '& .MuiOutlinedInput-root': {
                    backgroundColor: 'white',
                  }
                }}
              />

              <TextField
                fullWidth
                label="Пароль"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                margin="normal"
                required
                sx={{
                  '& .MuiOutlinedInput-root': {
                    backgroundColor: 'white',
                  }
                }}
              />

              <Button
                fullWidth
                type="submit"
                variant="contained"
                size="large"
                disabled={loading}
                startIcon={<LoginIcon />}
                sx={{
                  mt: 3,
                  mb: 2,
                  py: 1.5,
                  backgroundColor: '#3498db',
                  '&:hover': {
                    backgroundColor: '#2980b9',
                  },
                }}
              >
                {loading ? 'Вход...' : 'Войти'}
              </Button>
            </Box>
          </CardContent>
        </Card>

        <Paper 
          elevation={1} 
          sx={{ 
            mt: 3, 
            p: 3, 
            borderRadius: 2,
            backgroundColor: 'white',
          }}
        >
          <Typography variant="h6" gutterBottom sx={{ fontWeight: 600, color: '#2c3e50' }}>
            Демо-аккаунты:
          </Typography>
          <Divider sx={{ mb: 2 }} />
          <List dense>
            <ListItem sx={{ backgroundColor: '#f8f9fa', mb: 1, borderRadius: 1 }}>
              <ListItemText
                primary="Пациент"
                secondary="patient / password"
                primaryTypographyProps={{ fontWeight: 500, color: '#2c3e50' }}
                secondaryTypographyProps={{ color: '#7f8c8d' }}
              />
            </ListItem>
            <ListItem sx={{ backgroundColor: '#f8f9fa', mb: 1, borderRadius: 1 }}>
              <ListItemText
                primary="Клиника 1"
                secondary="clinic1 / password"
                primaryTypographyProps={{ fontWeight: 500, color: '#2c3e50' }}
                secondaryTypographyProps={{ color: '#7f8c8d' }}
              />
            </ListItem>
            <ListItem sx={{ backgroundColor: '#f8f9fa', borderRadius: 1 }}>
              <ListItemText
                primary="Регулятор"
                secondary="regulator / password"
                primaryTypographyProps={{ fontWeight: 500, color: '#2c3e50' }}
                secondaryTypographyProps={{ color: '#7f8c8d' }}
              />
            </ListItem>
          </List>
        </Paper>
      </Container>
    </Box>
  );
}

export default Login;
