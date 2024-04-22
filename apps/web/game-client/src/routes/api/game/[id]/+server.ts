import { error } from '@sveltejs/kit';
import { initNewGame, servers } from '$lib/server/sample-data/server-instance';
import crypto from 'crypto';
import type { RequestHandler } from './$types';

type PostAnswerPayload = {
  userId: string;
  gameId: number;
  answer: string;
}
export const POST: RequestHandler = async (event) => {
  const url = new URL(event.request.url);
  const pathSegments = url.pathname.split('/');
  const strId = pathSegments[pathSegments.length - 1];
  if (!strId) {
    return error(404);
  }
  const id = parseInt(strId);

  const server = servers.find((server) => server.id === id);
  if (!server) {
    return error(404, 'Server not found');
  }
  const game = server.game;

  const { gameId, userId, answer }: PostAnswerPayload = await event.request.json();
  if (gameId > server.gameData.length || gameId < 0) {
    return error(400, 'Invalid game ID');
  }
  if (!userId) {
    return error(401, 'Unauthorized');
  }
  server.users.add(userId);
  const headers = { 'Content-Type': 'application/json' };

  if (gameId !== game.id) {
    const gameResult = server.gameData[gameId];
    return new Response(JSON.stringify(gameResult), { headers });
  }
  const now = Date.now();

  game.logs.push({
    user: userId,
    answer: answer,
    time: now,
  })

  if (answer !== game.answer) {
    const body = JSON.stringify({
      id: game.id,
      question: game.question,
      rewards: game.logs.length,
      nextAt: game.nextAt,
    });
    return new Response(body, { headers });
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const {logs: _, ...gameData} = {
    ...game,
    winner: userId,
    rewards: game.logs.length,
    numOfTries: game.logs.reduce((acc, log) => acc + (log.user === userId ? 1 : 0), 0),
    finishedAt: now,
  }
  server.gameData.push(gameData)
  server.users = new Set()
  game.logs.forEach(log => server.users.add(log.user))

  const newGame = initNewGame();
  newGame.id = server.gameData.length;
  const hashString = crypto.createHash('sha256').update(newGame.startAt.toString()).digest('hex')
  newGame.question = hashString.slice(0, 8);
  newGame.answer = hashString.slice(hashString.length - 2);
  server.game = newGame;

  const gameResult = gameData;
  return new Response(JSON.stringify(gameResult), { headers });
};

// export const PUT: RequestHandler = (event) => {
//   // do something
//   const body = JSON.stringify({ message: 'OK' });
//   const headers = { 'Content-Type': 'application/json' };
//   return new Response(body, { headers });
// }

// export const PATCH: RequestHandler = (event) => {
//   // do something
//   const body = JSON.stringify({ message: 'OK' });
//   const headers = { 'Content-Type': 'application/json' };
//   return new Response(body, { headers });
// }

export const GET: RequestHandler = (event) => {
  const contentType = event.request.headers.get('Content-Type');
  if (contentType !== 'application/json') {
    return error(404)
  }

  const url = new URL(event.request.url);
  const pathSegments = url.pathname.split('/');
  const strId = pathSegments[pathSegments.length - 1];
  if (!strId) {
    return error(404);
  }
  const id = parseInt(strId);

  const server = servers.find((server) => server.id === id);
  if (!server) {
    return error(404, 'Server not found');
  }
  const params = new URLSearchParams(url.search);
  const paging = {
    offset: parseInt(params.get('offset') || '0') || 0,
    limit: parseInt(params.get('limit') || '10') || 10,
  }

  const gameData = server.gameData.slice().reverse().slice(paging.offset, paging.offset + paging.limit);
  const body = JSON.stringify({
    items: gameData,
    total: server.gameData.length,
    offset: paging.offset,
    limit: paging.limit,
  });
  const headers = { 'Content-Type': 'application/json' };
  return new Response(body, { headers });
}

