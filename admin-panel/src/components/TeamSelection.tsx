import { FormControl, InputLabel, MenuItem, Select, Stack, Typography, Grid, Switch, FormControlLabel, SelectChangeEvent, TextField, Autocomplete } from '@mui/material';
import { Team, Role, Player, GamePlayer } from '../types';

interface TeamSelectionProps {
  team: Team;
  players: Player[];
  selectedPlayers: GamePlayer[];
  onPlayerChange: (index: number, playerId: string | null, customNickname?: string) => void;
  onRoleChange: (index: number, role: Role) => void;
  onCaptainChange: (index: number) => void;
}

const ROLES: Role[] = ['carry', 'mid', 'offlane', 'pos4', 'pos5'];

export function TeamSelection({
  team,
  players,
  selectedPlayers,
  onPlayerChange,
  onRoleChange,
  onCaptainChange,
}: TeamSelectionProps) {
  return (
    <Stack spacing={2}>
      <Typography variant="h5" fontWeight="bold">
        {team} Team
      </Typography>
      {ROLES.map((role, index) => {
        const selectedPlayer = selectedPlayers[index];

        return (
          <Grid container key={role} spacing={2} alignItems="center">
            <Grid item xs={6}>
              <Autocomplete<Player, false, true, true>
                options={players}
                getOptionLabel={(option) => {
                  if (typeof option === 'string') return option;
                  return option.nickname;
                }}
                value={selectedPlayer?.player || null}
                isOptionEqualToValue={(option, value) => {
                  if (!option || !value) return false;
                  return option.id === value.id;
                }}
                onChange={(_, newValue) => {
                  if (newValue && typeof newValue !== 'string') {
                    onPlayerChange(index, newValue.id);
                  } else {
                    onPlayerChange(index, null);
                  }
                }}
                freeSolo
                onInputChange={(_, value) => {
                  if (value) {
                    const existingPlayer = players.find(
                      (p) => p.nickname.toLowerCase() === value.toLowerCase()
                    );
                    if (existingPlayer) {
                      onPlayerChange(index, existingPlayer.id);
                    } else {
                      onPlayerChange(index, 'custom', value);
                    }
                  }
                }}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label="Select or enter player"
                    fullWidth
                  />
                )}
              />
            </Grid>
            <Grid item xs={4}>
              <FormControl fullWidth>
                <InputLabel>Role</InputLabel>
                <Select
                  label="Role"
                  value={role}
                  onChange={(e) => onRoleChange(index, e.target.value as Role)}
                >
                  {ROLES.map((r) => (
                    <MenuItem key={r} value={r}>
                      {r}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={2}>
              <FormControlLabel
                control={
                  <Switch
                    checked={selectedPlayer?.isCaptain || false}
                    onChange={() => onCaptainChange(index)}
                  />
                }
                label="Captain"
              />
            </Grid>
          </Grid>
        );
      })}
    </Stack>
  );
} 