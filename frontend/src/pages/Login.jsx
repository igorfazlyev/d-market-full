import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useAuth } from '../context/AuthContext';
import { ROLES } from '../utils/constants';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    const result = await login(username, password);
    setLoading(false);

    if (result.success) {
      // Redirect based on role
      const { role } = result.user;
      switch (role) {
        case ROLES.PATIENT:
          navigate('/patient/dashboard');
          break;
        case ROLES.CLINIC:
          navigate('/clinic/dashboard');
          break;
        case ROLES.REGULATOR:
          navigate('/regulator/dashboard');
          break;
        default:
          navigate('/');
      }
    } else {
      setError(result.error);
    }
  };

  const handleDemoLogin = (demoUsername) => {
    setUsername(demoUsername);
    setPassword('password');
  };

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
        <Paper elevation={3} sx={{ p: 4, width: '100%' }}>
          <Typography variant="h4" align="center" gutterBottom>
            🦷 Стоматологический Маркетплейс
          </Typography>
          <Typography variant="subtitle1" align="center" color="text.secondary" sx={{ mb: 3 }}>
            Вход в систему
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Имя пользователя"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              margin="normal"
              required
              autoFocus
            />
            <TextField
              fullWidth
              label="Пароль"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              margin="normal"
              required
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ mt: 3, mb: 2 }}
            >
              {loading ? <CircularProgress size={24} /> : 'Войти'}
            </Button>
          </form>

          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle2" color="text.secondary" gutterBottom>
              Демо-аккаунты:
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Button
                variant="outlined"
                size="small"
                onClick={() => handleDemoLogin('patient')}
              >
                👤 Пациент
              </Button>
              <Button
                variant="outlined"
                size="small"
                onClick={() => handleDemoLogin('clinic1')}
              >
                🏥 Клиника 1
              </Button>
              <Button
                variant="outlined"
                size="small"
                onClick={() => handleDemoLogin('regulator')}
              >
                📊 Регулятор
              </Button>
            </Box>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default Login;
