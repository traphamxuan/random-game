<script lang="ts">
	import { GameStatus, storeGameState } from '$lib/client/stores/game-state';
    import type { PaginationSettings } from '@skeletonlabs/skeleton'
    import { tableMapperValues, Paginator, Table } from "@skeletonlabs/skeleton";
	import { onDestroy, onMount } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';

    export let serverID: string;

    let sourceData: (Omit<GameData, 'logs'> & { rewardStr: string })[] = [];
    let unsub: Unsubscriber;

    let paginationSettings = {
        page: 0,
        limit: 5,
        size: sourceData.length,
        amounts: [1,2,5,10],
    } satisfies PaginationSettings;

    onMount(async () => {
        const res = await fetch(`/api/game/${serverID}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!res.ok) {
            throw new Error(res.statusText)
        }
        const data = await res.json();
        paginationSettings.limit = data.limit;
        paginationSettings.size = data.total;
        paginationSettings.page = data.offset;
        sourceData = data.items.map((src: GameData) => ({
            ...src,
            rewardStr: `${src.rewards}/${src.numOfTries}`
        }));

        unsub = storeGameState.subscribe(val => {
            const gameData = $storeGameState.gameData;
            const updatedData = { ...gameData, rewardStr: `${gameData.rewards}` };
            if ($storeGameState.status === GameStatus.Finished) {
                updatedData.rewardStr += `/${gameData.numOfTries}`;
            }

            if (gameData.id !== sourceData[0]?.id) {
                sourceData.unshift(updatedData)
            } else {
                sourceData[0] = updatedData;
            }
            sourceData = sourceData;
        })

        paginationSettings = paginationSettings;
        const gameData = $storeGameState.gameData;
        if (gameData.id !== sourceData[0].id) {
            sourceData.unshift({
                id: gameData.id,
                question: gameData.question,
                answer: gameData.answer,
                winner: gameData.winner,
                numOfTries: gameData.numOfTries,
                rewards: gameData.rewards,
                finishedAt: gameData.finishedAt,
                startAt: gameData.startAt,
                rewardStr: `${gameData.rewards}`,
            })
        }
    });

    onDestroy(() => {
        if (unsub) {
            unsub();
        }
    })
</script>

<Table source={{
  // A list of heading labels.
  head: ['ID', 'Question', 'Answer', 'Winner', 'Reward'],
  // The data visibly shown in your table body UI.
  body: tableMapperValues(sourceData, ['id', 'question', 'answer', 'winner', 'rewardStr']),
  // Optional: The data returned when interactive is enabled and a row is clicked.
  // meta: tableMapperValues(sourceData, []),
  // Optional: A list of footer labels.
//   foot: ['Number of games', '', '<code class="code">5/123</code>']
}} />
<Paginator
    bind:settings={paginationSettings}
    showFirstLastButtons="{true}"
    showPreviousNextButtons="{true}"
/>
