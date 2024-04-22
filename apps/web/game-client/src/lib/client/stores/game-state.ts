import { writable } from "svelte/store";

export enum GameStatus {
  Waiting = "Waiting",
  Playing = "Playing",
  Finished = "Finished",
}

type GameState = {
  gameData: Omit<GameData, 'logs'>,
  status: GameStatus;
  answers: string[];
  numOfTries: number;
}

function initNewGame(): GameState {
  return {
    gameData: {
      id: 0,
      question: "",
      answer: "",
      rewards: 0,
      numOfTries: 0,
      startAt: 0,
      finishedAt: 0,
      winner: "",
    },
    status: GameStatus.Waiting,
    answers: [],
    numOfTries: 0,
  };
}

function createGameStateStore() {
  const { subscribe, set, update } = writable<GameState>(initNewGame());

  return {
    subscribe,
    start: (gameData: Pick<GameData, 'id' | 'question'>) => {
      const newGame = initNewGame();
      newGame.gameData.id = gameData.id;
      newGame.gameData.question = gameData.question;
      newGame.status = GameStatus.Playing;
      set(newGame);
    },
    run: (rewards: number) => update((state) => {
      state.gameData.rewards = rewards;
      return state;
    }),
    stop: (gameData: Omit<GameData, 'logs'>) => update((state) => {
      state.gameData = gameData;
      state.status = GameStatus.Finished;
      return state;
    }),
  };
}

export const storeGameState = createGameStateStore();