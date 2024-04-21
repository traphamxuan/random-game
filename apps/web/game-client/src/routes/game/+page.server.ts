import { servers } from '$lib/server/sample-data/server-instance';
import type { PageServerLoad } from './$types';

export const load = (async () => {
    return {
        servers: servers.map(server => ({
            id: server,
            name: server.name,
            capacity: server.capacity,
            size: server.users.size,
            createdAt: server.createdAt,
        })),
    };
}) satisfies PageServerLoad;
