<script lang="ts">
	import { type TableSource, tableMapperValues, Table } from "@skeletonlabs/skeleton";
	import type { PageData } from "./$types";
	import { goto } from "$app/navigation";

	export let data: PageData

	const sourceData = data.servers?.map((server, i) => {
		return {
			id: server.id,
			name: server.name,
			slots: `${server.size}/${server.capacity}`,
			lifetime: new Date().getTime() - server.createdAt,
			position: i+1,
		};
	});

	const tableSimple: TableSource = {
		// A list of heading labels.
		head: ['Name', 'Slots', 'Lifetime'],
		// The data visibly shown in your table body UI.
		body: tableMapperValues(sourceData, ['name', 'slots', 'lifetime']),
		// Optional: The data returned when interactive is enabled and a row is clicked.
		meta: tableMapperValues(sourceData, ['id', 'position', 'name']),
		// Optional: A list of footer labels.
		foot: ['Total', '', '<code class="code">5</code>']
	};

	const gotoServer = (event: CustomEvent) => {
		console.log(event)
		goto(`/game/${event.detail[1]}`);
	}
				
</script>

<pre class="pre">Select server to start the game</pre>
<Table source={tableSimple} interactive on:selected={gotoServer}/>