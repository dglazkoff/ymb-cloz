export type Role = 'carry' | 'mid' | 'offlane' | 'pos4' | 'pos5';
export type Team = 'RADIANT' | 'DIRE';

export interface Player {
  id: string;
  nickname: string;
  games?: number;
}

export interface GamePlayer {
  player: Player;
  team: Team;
  role: Role;
  isCaptain: boolean;
}

export interface Game {
  id?: string;
  timestamp: Date;
  radiantPlayers: GamePlayer[];
  direPlayers: GamePlayer[];
  winningTeam: Team;
} 