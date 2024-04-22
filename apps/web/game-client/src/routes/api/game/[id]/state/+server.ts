import { error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { servers } from "$lib/server/sample-data/server-instance";

export const GET: RequestHandler = (event) => {
  const url = new URL(event.request.url);
  const pathSegments = url.pathname.split('/');
  const strId = pathSegments[pathSegments.length - 2];
  if (!strId) {
    return error(404);
  }
  const id = parseInt(strId);

  const server = servers.find((server) => server.id === id);
  if (!server) {
    return error(404, 'Server not found');
  }
  const game = server.game;
  const body = JSON.stringify({
    id: game.id,
    question: game.question,
    rewards: game.logs.length,
    nextAt: game.nextAt,
  });
  console.log(game.question, game.answer)
  const headers = { 'Content-Type': 'application/json' };
  return new Response(body, { headers });
}
