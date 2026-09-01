<script lang="ts">
  import { page } from '$app/state'
  import type { stageTime } from '$lib/components/stage-time/stage-time-row.svelte';
  import { getDurationMsFromStageTime, isStageTimeRunning } from '$lib/components/stage-time/stage-time-row.svelte';
  import * as StageTime from "$lib/components/stage-time"
  import * as Tabs from "$lib/components/ui/tabs"
  import { pb } from '$lib/pb'
  import type { RecordModel } from 'pocketbase'
  import { onMount } from 'svelte'
  import { toast } from 'svelte-sonner'
  import { SvelteMap } from 'svelte/reactivity';

  let stages: RecordModel[] = $state([])
  let currentStageTab = $state('')

  let stageTimes = $state(new SvelteMap<string /* Transponder */, stageTime>)
  let currentStageTimesStage = ''

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

  $effect(() => {
    if (currentStageTab != '') {
      if (currentStageTab != currentStageTimesStage) {
        currentStageTimesStage = currentStageTab
        requestStageTimes()
      }
    }
  })

  function requestStageTimes () {
    pb.collection("stage_times").unsubscribe()
    stageTimes.clear()

    const filter = currentStageTab == "all" ?
      pb.filter("stage.event = {:event}", {"event": page.params.eventId}) :
      pb.filter("stage = {:stage}", {"stage": currentStageTab})

    pb.collection("stage_times").getFullList({
      filter: filter,
      expand: "start, stop, penalties, start.transponder.owner"
    })
      .then((result) => {
        parseStageTimes(result, "initialLoad")
      })
      .catch(() => {
        toast.error("Couldn't request Stage Times")
      })
    
    pb.collection("stage_times").subscribe("*", (newRecord) => {
      parseStageTimes([newRecord.record], newRecord.action)
    }, {
      filter: filter,
      expand: "start, stop, penalties, start.transponder.owner"
    })
      .catch(() => {
        toast.error("Couldn't subscribe to Stage Times")
      })
  }

  function parseStageTimes (newStageTimes: RecordModel[], action: string) {
    for (const newStageTime of newStageTimes) {
      const transponder: string = newStageTime.expand?.start.transponder

      if (action == "delete") {
        deleteRecordFromStageTime(transponder, newStageTime)
        continue
      }

      if (!stageTimes.has(transponder)) {
        if (newStageTime.stop == '') {
          stageTimes.set(transponder, {running: newStageTime} as stageTime)
        } else {
          stageTimes.set(transponder, {
            last: newStageTime,
            best: newStageTime
          } as stageTime)
        }
        continue
      }

      const existingStageTimes = stageTimes.get(transponder) as stageTime

      if (isStageTimeRunning(newStageTime)) {
        existingStageTimes.running = newStageTime
      } else {
        existingStageTimes.last = newStageTime

        if (existingStageTimes.running?.id == newStageTime.id) {
          existingStageTimes.running = undefined
        }

        if (existingStageTimes.best == undefined) {
          existingStageTimes.best = newStageTime
        } else {
          if (getDurationMsFromStageTime(newStageTime) < getDurationMsFromStageTime(existingStageTimes.best as RecordModel)) {
            existingStageTimes.best = newStageTime
          }
        }
      }
    }

    sortStageTimesByBest()
  }

  function sortStageTimesByBest () {
    stageTimes = new SvelteMap(Array.from(stageTimes).sort((a, b) => {
      if (a[1].best != undefined && b[1].best != undefined) {
        return getDurationMsFromStageTime(a[1].best).valueOf() - getDurationMsFromStageTime(b[1].best).valueOf()
      } else if (a[1].best == undefined && b[1].best == undefined) {
        return 0
      } else if (a[1].best != undefined) {
        return -1
      } else {
        return 1
      }
    }))
  }

  function deleteRecordFromStageTime (transponder: string, recordToDelete: RecordModel) {
    const existingStageTime = stageTimes.get(transponder)

    if (existingStageTime == undefined) {
      return
    }

    if (existingStageTime.best?.id == recordToDelete.id) {
      existingStageTime.best = undefined
      // TODO Request new best
    }

    if (existingStageTime.last?.id == recordToDelete.id) {
      existingStageTime.last = undefined
      // TODO Request new last
    }

    if (existingStageTime.running?.id == recordToDelete.id) {
      existingStageTime.running = undefined
    }

    if (
      existingStageTime.best == undefined &&
      existingStageTime.last == undefined &&
      existingStageTime.running == undefined
    ) {
      stageTimes.delete(transponder)  
    }
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
        <!-- <Tabs.Trigger value="all">All</Tabs.Trigger> --> <!-- TODO Needs different Stage Time Table variant -->
      </Tabs.List>
    {/if}
  </Tabs.Root>
  
  <StageTime.Root>
    {#each stageTimes, index}
      <StageTime.Row position={index + 1} stageTimes={stageTimes} />
    {/each}
  </StageTime.Root>

</div>