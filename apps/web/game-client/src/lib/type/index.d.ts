type GameAction = {
  user: string;
  answer: string;
  time: number;
}

type GameData = {
  id: number;
  question: string;
  answer: string;
  winner: string;
  numOfTries: number;
  rewards: number;
  startAt: number;
  finishedAt: number;
  logs: GameAction[];
}

type GameRunning = {
  id: number;
  question: string;
  answer: string;
  rewards: number;
  nextAt: number;
  startAt: number;
  logs: GameAction[];
}

type ServerData = {
  id: number;
  name: string;
  capacity: number;
  createdAt: number;
  game: GameRunning;
  gameData: GameData[];
  users: Set<string>;
};
