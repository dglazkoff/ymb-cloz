import { Player } from '../types';
import { API_BASE_URL } from './gameService';

class PlayerService {
  private readonly baseUrl = API_BASE_URL;

  async getPlayers(): Promise<{ players: Player[] }> {
    const response = await fetch(`${this.baseUrl}/players`);
    if (!response.ok) {
      throw new Error('Failed to fetch players');
    }
    return response.json();
  }
}

export const playerService = new PlayerService(); 