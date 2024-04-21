<script lang="ts">
	import { page } from '$app/stores';
    import InputPlayer from './InputPlayer.svelte';
	import TableResults from './TableResults.svelte';

</script>

<InputPlayer serverID={$page.params.slug}/>
<TableResults serverID={$page.params.slug} />