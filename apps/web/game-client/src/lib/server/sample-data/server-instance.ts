
export const initNewGame = (): GameRunning => ({
  id: 0,
  question: 'aa',
  answer: '96',
  nextAt: 0,
  startAt: Date.now(),
  logs: [],
});

export const servers: ServerData[] = [
  { id: 1, name: 'Hydrogen', capacity:50, createdAt: 1713690933000, gameData: [], users: new Set(), game: initNewGame() },
  { id: 2, name: 'Helium', capacity: 50, createdAt: 1713690933000, gameData: [], users: new Set(), game: initNewGame() },
  { id: 3, name: 'Lithium', capacity: 50, createdAt: 1713690933000, gameData: [], users: new Set(), game: initNewGame() },
  { id: 4, name: 'Beryllium', capacity: 50, createdAt: 1713690933000, gameData: [], users: new Set(), game: initNewGame() },
  { id: 5, name: 'Boron', capacity:50, createdAt: 1713690933000, gameData: [], users: new Set(), game: initNewGame() },
];
