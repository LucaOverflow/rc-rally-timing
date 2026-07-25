<script lang="ts">
  import { page } from '$app/state'
  import * as Tabs from "$lib/components/ui/tabs"
  import { pb } from '$lib/pb'
  import type { RecordModel } from 'pocketbase'
  import { onMount } from 'svelte'
  import { toast } from 'svelte-sonner'

  let stages: RecordModel[] = $state([])
  let currentStageTab = $state('')

  onMount(() => {
    requestStages()
  })

  function requestStages () {
    pb.collection('stages').getFullList({
      filter: pb.filter('event = {:eventId}', {eventId: page.params.eventId})
    })
      .then((result: RecordModel[]) => {
        stages = result
        defaultSelectStage()
      })
      .catch(() => {
        toast.error("Couldn't load stages")
      })
  }

  function defaultSelectStage () {
    if (currentStageTab != '') {
      return
    }

    // 1. Select first active
    for (const stage of stages) {
      if (stage.active) {
        currentStageTab = stage.id
        return
      }
    }

    // 2. Select all
    currentStageTab = 'all'
  }
</script>

<div class="ml-6 mr-6">
  <Tabs.Root bind:value={currentStageTab}>
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
        <Tabs.Trigger value="all">All</Tabs.Trigger>
      </Tabs.List>
    {/if}
    {#each stages as stage}
      <Tabs.Content value={stage.id}>

        {stage.name}

      </Tabs.Content>
    {/each}
    <Tabs.Content value="all">

      All

    </Tabs.Content>
  </Tabs.Root>
</div>