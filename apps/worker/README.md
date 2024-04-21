This worker in the background waiting for request trigger via grpc (upgrade to kafka later). Any request via grpc will trigger this server to collect data in high-write game DB, then calculate and publish the leader board of server to the leader board database of game-api.

