<script lang="ts">
  import { page } from '$app/state'
  import * as Tabs from "$lib/components/ui/tabs"
  import { pb } from '$lib/pb'
  import type { RecordModel } from 'pocketbase'
  import { onMount } from 'svelte'
  import { toast } from 'svelte-sonner'

  let stages: RecordModel[] = $state([])
  let currentStageTab = $state('')

  interface stageTime {
    best: RecordModel | undefined,
    last: RecordModel | undefined,
    running: RecordModel | undefined
  }
  let stageTimes = $state(new Map<string /* Transponder */, stageTime>)
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
      expand: "start, stop, penalties"
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
      expand: "start, stop, penalties"
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
          if (getDurationFromStageTime(newStageTime) < getDurationFromStageTime(existingStageTimes.best as RecordModel)) {
            existingStageTimes.best = newStageTime
          }
        }
      }
    }
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

  function getDurationFromStageTime (stageTime: RecordModel): Date {
    let stopTime_ms: number
    if (stageTime.stop == "") {
      stopTime_ms = Date.now()
    } else {
      stopTime_ms = stageTime.expand?.stop.timecode_ms
    }
    
    let penalty_ms = 0
    if (stageTime.penalties.length > 0) {
      for (const penalty of stageTime.expand?.penalties) {
        penalty_ms += penalty.duration_ms
      }
    }

    return new Date((stopTime_ms - stageTime.expand?.start.timecode_ms) + penalty_ms)
  }

  function isStageTimeRunning (stageTime: RecordModel): boolean {
    return stageTime.stop == ""
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