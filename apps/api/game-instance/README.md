Game services support client playing the game. All endpoints must be authorized

Initially, this server generates a set of requirements randomly. These requirements is broadcasted for all clients who is connecting to this servers.

When a client submit an answer to this server, it stores the answer to high-write database and verify the answer. If answer is right, server will trigger worker to collect the results and refresh requirements.

Internal request via grpc from game-api to register token before access the game
0. Register client (/api/game/:server_id/register)

All requests below must be verified if the token to access the game is valid
1. Start game at specific server (/api/game/:server_id)
2. Fetch challenge requirements (/api/game/:server_id/requirements)
3. Submit result (/api/game/:server_id/result)
4. Fetch status (/api/game/:server_id/status)
