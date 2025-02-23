import { CreateGameRequest } from '../types/game';

export const API_BASE_URL = 'http://localhost:8080/api';

export const gameService = {
    createGame: async (gameData: CreateGameRequest): Promise<void> => {
        const response = await fetch(`${API_BASE_URL}/games`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(gameData),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.json();
    },
}; 