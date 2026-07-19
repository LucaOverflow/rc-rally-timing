<script lang="ts">
  import { page } from '$app/state';
  import { Button } from "$lib/components/ui/button"
  import * as Tabs from "$lib/components/ui/tabs"
  import * as Dialog from "$lib/components/ui/dialog"
  import { Label } from "$lib/components/ui/label"
  import { Input } from "$lib/components/ui/input"
  import { pb } from '$lib/pb';
  import { PlusSignIcon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import type { RecordModel } from 'pocketbase';
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';

  let stages: RecordModel[] = $state([])
  let currentStageTab = $state('')
  let openAddStagePopup = $state(false)
  let addStageFormData = $state({
    name: ''
  })

  onMount(() => {
    requestStages()
  })

  function requestStages () {
    pb.collection('stages').getFullList({
      filter: pb.filter('event = {:eventId}', {eventId: page.params.eventId})
    })
      .then((result: RecordModel[]) => {
        stages = result
      })
      .catch(() => {
        toast.error("Couldn't load stages")
      })
  }

  function addStage () {
    pb.collection('stages').create({
      name: addStageFormData.name,
      event: page.params.eventId
    })
      .then((result) => {
        requestStages()
        currentStageTab = result.id
        openAddStagePopup = false
        addStageFormData.name = ''
      })
      .catch(() => {
        toast.error("Couldn't add stage")
      })
  }
</script>

<div class="ml-6 mr-6">
  <Tabs.Root bind:value={currentStageTab}>
    <div class="flex gap-2">
      <Button onclick={() => {openAddStagePopup = true}}>
        <HugeiconsIcon icon={PlusSignIcon} />
        Add Stage
      </Button>
      {#if stages.length > 0}
        <Tabs.List class="ml-auto overflow-x-auto scrollbar-none">
          {#each stages as stage}
            <Tabs.Trigger value={stage.id}>
              {#if stage.active}
                <div class="w-2 h-2 rounded-full" style="background-color: red"></div>
              {/if}
              {stage.name}
            </Tabs.Trigger>
          {/each}
        </Tabs.List>
      {/if}
    </div>
    {#each stages as stage}
      <Tabs.Content value={stage.id}>
        <p>This is the content of {stage.name}</p>
      </Tabs.Content>
    {/each}
  </Tabs.Root>
</div>

<!-- Add Stage Popup -->
<Dialog.Root bind:open={openAddStagePopup}>
  <Dialog.Content>
    <Dialog.Title>Add Stage</Dialog.Title>
    <form onsubmit={addStage}>
      <Dialog.Header>

        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={addStageFormData.name} />

      </Dialog.Header>
      <Dialog.Footer>
        <Dialog.Close>Cancel</Dialog.Close>
        <Button type="submit">Add</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>