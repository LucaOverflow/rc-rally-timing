<script lang="ts">
  import { goto } from '$app/navigation'
  import { Button } from "$lib/components/ui/button"
  import * as Item from '$lib/components/ui/item'
  import * as Dialog from "$lib/components/ui/dialog"
  import * as Alert from "$lib/components/ui/alert"
  import { ScrollArea } from '$lib/components/ui/scroll-area'
  import { globals } from '$lib/globals.svelte'
  import { pb } from '$lib/pb'
  import { Alert01Icon, PlusSignIcon, StopWatchIcon } from '@hugeicons/core-free-icons'
  import { HugeiconsIcon } from '@hugeicons/svelte'
  import type { ClientResponseError, ListResult, RecordModel } from 'pocketbase'
  import { onMount } from 'svelte'
  import { toast } from 'svelte-sonner'
  import EventSettings from '$lib/components/event-settings.svelte';

  let events: RecordModel[] = $state([])

  const emptyAddEventFormData = {
    name: "",
    description: "",
    start: "",
    end: "",
    location: "",
    start_order: "",
    start_mode: "",
    signup_until: "",
    max_drivers: 0,
    allow_transponder_signup: true
  }

  let openAddEventPopup = $state(false)
  let addEventFormData = $state(emptyAddEventFormData)
  let addEventErrorMessage = $state("")

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

  function addEvent () {
    addEventErrorMessage = ""

    pb.collection("events").create({
      organizer: [pb.authStore.record?.id],
      ...addEventFormData
    })
      .then(() => {
        openAddEventPopup = false
        addEventFormData = emptyAddEventFormData
        requestEvents()
      })
      .catch((error: ClientResponseError) => {
        addEventErrorMessage = error.message
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

<Button onclick={() => {openAddEventPopup = true}} class="mt-2">
  <HugeiconsIcon icon={PlusSignIcon} />
  Add
</Button>

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

<!-- Add Event Popup -->
<Dialog.Root bind:open={openAddEventPopup}>
  <Dialog.Content class="flex max-h-[90%] flex-col gap-0">
    <Dialog.Header class="overflow-y-auto">
      
      <EventSettings title="Add a new Event" description="You can still edit all these settings after creation." bind:formData={addEventFormData} />

      {#if addEventErrorMessage != ''}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to add Event</Alert.Title>
          <Alert.Description>
            {addEventErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button onclick={addEvent}>Add</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>