<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import * as Sidebar from '$lib/components/ui/sidebar/index';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import { Toaster } from "$lib/components/ui/sonner"

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<Toaster position="bottom-center" />

<Sidebar.Provider>
	<AppSidebar />
	<main class="flex-grow min-w-0">
		<Sidebar.Trigger />
		{@render children()}
	</main>
</Sidebar.Provider>
