import type { PageServerLoad } from './$types';

export const load = (async () => {
    return {
        servers: [
            { id: 1, name: 'Hydrogen', size: 18, capacity:50, createdAt: 1713690933000 },
            { id: 2, name: 'Helium', size: 40, capacity: 50, createdAt: 1713690933000 },
            { id: 3, name: 'Lithium', size: 26, capacity: 50, createdAt: 1713690933000 },
            { id: 4, name: 'Beryllium', size: 19, capacity: 50, createdAt: 1713690933000 },
            { id: 5, name: 'Boron', size: 10, capacity:50, createdAt: 1713690933000 },
        ],
    };
}) satisfies PageServerLoad;