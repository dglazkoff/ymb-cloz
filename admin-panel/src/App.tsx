import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { GameForm } from './components/GameForm';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <GameForm />
    </ThemeProvider>
  );
}

export default App;
