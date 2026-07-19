<script lang="ts">
  import { page } from '$app/state';
  import { Button } from "$lib/components/ui/button"
  import * as Tabs from "$lib/components/ui/tabs"
  import * as Dialog from "$lib/components/ui/dialog"
  import { Label } from "$lib/components/ui/label"
  import { Input } from "$lib/components/ui/input"
  import { pb } from '$lib/pb';
  import { Delete01Icon, PlayIcon, PlusSignIcon, StopIcon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import type { RecordModel } from 'pocketbase';
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';
  import DecoderSelector from '$lib/components/decoder-selector.svelte';

  let stages: RecordModel[] = $state([])
  let currentStageTab = $state('')
  let openAddStagePopup = $state(false)
  let addStageFormData = $state({
    name: ''
  })

  let openStartStagePopup = $state(false)
  let startStageFormData = $state({
    id: '',
    start_decoder: '',
    stop_decoder: ''
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

  function addStage (e: Event) {
    e.preventDefault()

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

  function removeStage (stageId: string) {
    pb.collection('stages').delete(stageId)
      .then(() => {
        requestStages()
      })
      .catch(() => {
        toast.error("Couldn't remove stage")
      })
  }

  function startStage (e: Event) {
    e.preventDefault()

    pb.collection('stages').update(startStageFormData.id, {
      start_decoder: startStageFormData.start_decoder,
      stop_decoder: startStageFormData.stop_decoder,
      active: true
    })
      .then(() => {
        requestStages()
        openStartStagePopup = false
        startStageFormData.id = ''
        startStageFormData.start_decoder = ''
        startStageFormData.stop_decoder = ''
      })
      .catch(() => {
        toast.error("Couldn't start stage")
      })
  }

  function stopStage (stageId: string) {
    pb.collection('stages').update(stageId, {
      active: false,
      start_decoder: '',
      stop_decoder: ''
    })
      .then(() => {
        requestStages()
      })
      .catch(() => {
        toast.error("Couldn't stop stage")
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

        {#if stage.active}
          <Button variant="destructive" onclick={() => {stopStage(stage.id)}}>
            <HugeiconsIcon icon={StopIcon} />
            Stop
          </Button>
        {:else}
          <Button onclick={() => {startStageFormData.id = stage.id; openStartStagePopup = true}}>
            <HugeiconsIcon icon={PlayIcon} />
            Start
          </Button>
        {/if}

        <Button variant="destructive" onclick={() => {removeStage(stage.id)}}>
          <HugeiconsIcon icon={Delete01Icon} />
          Remove
        </Button>

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

<!-- Start Stage Popup -->
<Dialog.Root bind:open={openStartStagePopup}>
  <Dialog.Content>
    <Dialog.Title>Start Stage</Dialog.Title>
    <form onsubmit={startStage}>
      <Dialog.Header>

        <div class="grid grid-cols-2 items-baseline gap-2">
          <span>Start Decoder</span>
          <DecoderSelector bind:value={startStageFormData.start_decoder} />

          <span>Stop Decoder</span>
          <DecoderSelector bind:value={startStageFormData.stop_decoder} />
        </div>

      </Dialog.Header>
      <Dialog.Footer>
        <Dialog.Close>Cancel</Dialog.Close>
        <Button type="submit">Start</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>