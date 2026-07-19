<script lang="ts">
  import { goto } from '$app/navigation';
  import { Button } from "$lib/components/ui/button"
  import * as Item from '$lib/components/ui/item';
  import { ScrollArea } from '$lib/components/ui/scroll-area';
  import { globals } from '$lib/globals.svelte';
  import { pb } from '$lib/pb';
  import { StopWatchIcon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import type { ListResult, RecordModel } from 'pocketbase';
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';

  let events: RecordModel[] = $state([])

  onMount(() => {
    requestEvents()
  })

  function requestEvents () {
    pb.collection('events').getList(1, 25, {
      filter: pb.filter('organizer ~ {:user}', {user: pb.authStore.record?.id}),
      sort: '-created'
    })
      .then((result: ListResult<RecordModel>) => {
        events = result.items
      })
      .catch(() => {
        toast.error("Couldn't load Events")
      })
  }

  function quickAddPractice () {
    const date = new Date()

    pb.collection('events').create({
      name: "Quick practice " + date.toLocaleDateString(),
      start_order: "reactive",
      start_mode: "passing",
      organizer: pb.authStore.record?.id,
      allow_transponder_signup: true
    })
      .then(() => {
        requestEvents()
      })
      .catch(() => {
        toast.error("Couldn't add event")
      })
  }

  function deleteEvent (event: string) {
    pb.collection('events').delete(event)
      .then(() => {
        requestEvents()
        if (globals.activeEvent?.id == event) {
          globals.activeEvent = undefined
        }
      })
      .catch(() => {
        toast.error("Couldn't remove event")
      })
  }
</script>

<Button onclick={quickAddPractice} class="mt-2">
  <HugeiconsIcon icon={StopWatchIcon} />
  Quick Add Practice
</Button>

<ScrollArea class="m-6">
  {#each events as event, i}
    <Item.Root variant="outline" class="mb-2">
      <Item.Content class="flex-row space-x-1 items-center">
        <Item.Title>{event.name}</Item.Title>
        <div class="flex flex-col gap-1 ml-auto">
          <Button onclick={() => {goto('/events/' + event.id + '/schedule')}}>Select</Button>
          <Button variant="destructive" onclick={() => {deleteEvent(event.id)}}>Remove</Button>
        </div>
      </Item.Content>
    </Item.Root>
  {/each}
</ScrollArea>

<!-- TODO Add pagination navigation here -->