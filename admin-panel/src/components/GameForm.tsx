import { useState, useEffect } from 'react';
import { Button, Container, Stack, Paper, Typography, ToggleButtonGroup, ToggleButton } from '@mui/material';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { Team, Player, GamePlayer, Role } from '../types';
import { TeamSelection } from './TeamSelection';
import { gameService } from '../services/gameService';
import { playerService } from '../services/playerService';
import dayjs from 'dayjs';
import { CreateGameRequest, GamePlayerInput } from '../types/game';

export function GameForm() {
  const [timestamp, setTimestamp] = useState<Date>(new Date());
  const [winningTeam, setWinningTeam] = useState<Team>('RADIANT');
  const [radiantPlayers, setRadiantPlayers] = useState<GamePlayer[]>([]);
  const [direPlayers, setDirePlayers] = useState<GamePlayer[]>([]);
  const [players, setPlayers] = useState<Player[]>([]);

  useEffect(() => {
    const fetchPlayers = async () => {
      try {
        const fetchedPlayers = await playerService.getPlayers();
        setPlayers(fetchedPlayers.players ?? []);
      } catch (error) {
        console.error('Error fetching players:', error);
      }
    };
    fetchPlayers();
  }, []);

  const handlePlayerChange = (team: Team, index: number, playerId: string | null, customNickname?: string) => {
    if (!playerId) return;
    
    let player: Player;
    if (playerId === 'custom' && customNickname !== undefined) {
      player = {
        id: 'custom',
        nickname: customNickname
      };
    } else {
      const existingPlayer = players.find(p => p.id === playerId);
      if (!existingPlayer) return;
      player = existingPlayer;
    }

    const newGamePlayer: GamePlayer = {
      player,
      team,
      role: ['carry', 'mid', 'offlane', 'pos4', 'pos5'][index] as Role,
      isCaptain: false
    };

    if (team === 'RADIANT') {
      const newPlayers = [...radiantPlayers];
      newPlayers[index] = newGamePlayer;
      setRadiantPlayers(newPlayers);
    } else {
      const newPlayers = [...direPlayers];
      newPlayers[index] = newGamePlayer;
      setDirePlayers(newPlayers);
    }
  };

  const handleRoleChange = (team: Team, index: number, role: Role) => {
    if (team === 'RADIANT') {
      const newPlayers = [...radiantPlayers];
      if (newPlayers[index]) {
        newPlayers[index] = { ...newPlayers[index], role };
        setRadiantPlayers(newPlayers);
      }
    } else {
      const newPlayers = [...direPlayers];
      if (newPlayers[index]) {
        newPlayers[index] = { ...newPlayers[index], role };
        setDirePlayers(newPlayers);
      }
    }
  };

  const handleCaptainChange = (team: Team, index: number) => {
    if (team === 'RADIANT') {
      const newPlayers = radiantPlayers.map((p, i) => ({
        ...p,
        isCaptain: i === index
      }));
      setRadiantPlayers(newPlayers);
    } else {
      const newPlayers = direPlayers.map((p, i) => ({
        ...p,
        isCaptain: i === index
      }));
      setDirePlayers(newPlayers);
    }
  };

  const handleSubmit = async () => {
    // Transform the data to match the API format
    const transformPlayers = (players: GamePlayer[]): GamePlayerInput[] => players.map(p => ({
      ...(p.player.id === 'custom' ? { nickname: p.player.nickname } : { id: p.player.id }),
      role: p.role.toLowerCase(),
      is_captain: p.isCaptain
    }));

    const gameData: CreateGameRequest = {
      radiant_players: transformPlayers(radiantPlayers),
      dire_players: transformPlayers(direPlayers),
      winner: winningTeam
    };

    try {
      await gameService.createGame(gameData);

      // Reset form
      setRadiantPlayers([]);
      setDirePlayers([]);
      setTimestamp(new Date());
      setWinningTeam('RADIANT');
    } catch (error) {
      console.error('Error saving game:', error);
    }
  };

  return (
    <Container maxWidth="lg">
      <Paper elevation={1} sx={{ p: 3, borderRadius: 2 }}>
        <Stack spacing={4}>
          <Typography variant="h4">Record New Game</Typography>
          
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DateTimePicker
              label="Game Date and Time"
              value={dayjs(timestamp)}
              onChange={(date) => date && setTimestamp(date.toDate())}
            />
          </LocalizationProvider>

          <TeamSelection
            team="RADIANT"
            players={players}
            selectedPlayers={radiantPlayers}
            onPlayerChange={(index, playerId, customNickname) => handlePlayerChange('RADIANT', index, playerId, customNickname)}
            onRoleChange={(index, role) => handleRoleChange('RADIANT', index, role)}
            onCaptainChange={(index) => handleCaptainChange('RADIANT', index)}
          />

          <TeamSelection
            team="DIRE"
            players={players}
            selectedPlayers={direPlayers}
            onPlayerChange={(index, playerId, customNickname) => handlePlayerChange('DIRE', index, playerId, customNickname)}
            onRoleChange={(index, role) => handleRoleChange('DIRE', index, role)}
            onCaptainChange={(index) => handleCaptainChange('DIRE', index)}
          />

          <Stack spacing={2}>
            <Typography variant="h6">Winner</Typography>
            <ToggleButtonGroup
              value={winningTeam}
              exclusive
              onChange={(_, value) => value && setWinningTeam(value as Team)}
              fullWidth
            >
              <ToggleButton value="RADIANT">
                Radiant Victory
              </ToggleButton>
              <ToggleButton value="DIRE">
                Dire Victory
              </ToggleButton>
            </ToggleButtonGroup>
          </Stack>

          <Button variant="contained" size="large" onClick={handleSubmit} fullWidth>
            Save Game
          </Button>
        </Stack>
      </Paper>
    </Container>
  );
} 