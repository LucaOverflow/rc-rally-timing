<script lang="ts">
  import * as Empty from "$lib/components/ui/empty/index"
  import { pb } from '$lib/pb';
  import { RacingFlagIcon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import { onMount } from 'svelte';

  let isLoggedIn = $state(pb.authStore.isValid)

  onMount(() => {
    pb.authStore.onChange(() => {
      isLoggedIn = pb.authStore.isValid
    })
  })
</script>

<Empty.Root class="h-full flex items-center justify-center">
  <Empty.Header>
    <Empty.Media variant="icon">
      <HugeiconsIcon icon={RacingFlagIcon} />
    </Empty.Media>
    {#if isLoggedIn}
      <Empty.Title>Hey {pb.authStore.record?.name}</Empty.Title>
      <Empty.Description>
        Select a page from the menu to get started.
      </Empty.Description>
    {:else}
      <Empty.Title>Welcome</Empty.Title>
      <Empty.Description>
        With this tool you can manage and run RC Rally Events. Login or Register to get started.
      </Empty.Description>
    {/if}
  </Empty.Header>
</Empty.Root>
