<script lang="ts">
	import { type TableSource, tableMapperValues, Paginator, type PaginationSettings, Table } from "@skeletonlabs/skeleton";

	const sourceData = [
        {
            id: 1,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 2,
            question: 'f1',
            answer: 'a1',
            winner: 'Tra Pham',
            numOfTries: 13,
            rewards: 1024,
        },{
            id: 3,
            question: 'b2',
            answer: 'a6',
            winner: 'Tra Pham',
            numOfTries: 19,
            rewards: 1024,
        },{
            id: 4,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 244,
            rewards: 1024,
        },{
            id: 5,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 6,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 7,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 8,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 9,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },{
            id: 10,
            question: 'a3',
            answer: '1f',
            winner: 'Tra Pham',
            numOfTries: 14,
            rewards: 1024,
        },
    ];

    const sourceResults = sourceData.map((src) => ({
        ...src,
        rewardStr: `${src.rewards}/${src.numOfTries}`,
    }))

    let paginationSettings = {
        page: 0,
        limit: 5,
        size: sourceResults.length,
        amounts: [1,2,5,10],
    } satisfies PaginationSettings;

    let paginatedSource: (typeof sourceResults) = [];
    $: paginatedSource = sourceResults.slice(
        paginationSettings.page * paginationSettings.limit,
        paginationSettings.page * paginationSettings.limit + paginationSettings.limit
    );


    const tableSimple: TableSource = {
		// A list of heading labels.
		head: ['ID', 'Question', 'Answer', 'Winner', 'Reward'],
		// The data visibly shown in your table body UI.
		body: tableMapperValues(paginatedSource, ['id', 'question', 'answer', 'winner', 'rewardStr']),
		// Optional: The data returned when interactive is enabled and a row is clicked.
		// meta: tableMapperValues(sourceData, []),
		// Optional: A list of footer labels.
		foot: ['Number of games', '', '<code class="code">5/123</code>']
	};

</script>

<Table source={tableSimple} />
<Paginator
	bind:settings={paginationSettings}
	showFirstLastButtons="{true}"
	showPreviousNextButtons="{true}"
/>
