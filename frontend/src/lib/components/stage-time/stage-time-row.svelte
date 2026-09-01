<script module>
  export interface stageTime {
    best: RecordModel | undefined,
    last: RecordModel | undefined,
    running: RecordModel | undefined
  }

  export function getDurationMsFromStageTime (stageTime: RecordModel): number {
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

    return (stopTime_ms - stageTime.expand?.start.timecode_ms) + penalty_ms
  }

  export function isStageTimeRunning (stageTime: RecordModel): boolean {
    return stageTime.stop == ""
  }
</script>

<script lang="ts">
  import type { RecordModel } from 'pocketbase'
  import { Badge } from '$lib/components/ui/badge'
  import * as Table from "$lib/components/ui/table"
  import type { SvelteMap } from 'svelte/reactivity';

  let { position, stageTimes }: {position: number, stageTimes: SvelteMap<string /* Transponder */, stageTime>} = $props()
  let stageTime = $derived(stageTimes.values().toArray()[position-1])

  function getDriverName (): string | undefined {
    if (stageTime.last != undefined) {
      if (stageTime.last.expand?.start?.expand?.transponder?.id) {
        if (stageTime.last.expand?.start?.expand?.transponder?.expand?.owner?.name) {
          return stageTime.last.expand.start.expand.transponder.expand.owner.name
        } else {
          return stageTime.last.expand.start.expand.transponder.id
        }
      }
    } else if (stageTime.running != undefined) {
      if (stageTime.running.expand?.start?.expand?.transponder?.id) {
        if (stageTime.running.expand?.start?.expand?.transponder?.expand?.owner?.name) {
          return stageTime.running.expand.start.expand.transponder.expand.owner.name
        } else {
          return stageTime.running.expand.start.expand.transponder.id
        }
      }
    }
  }

  function getDifferentiator (): string | undefined {
    if (stageTime.last != undefined) {
      if (stageTime.last.expand?.start?.expand?.transponder?.id) {
        return stageTime.last.expand.start.expand.transponder.differentiator
      }
    } else if (stageTime.running != undefined) {
      if (stageTime.running.expand?.start?.expand?.transponder?.id) {
        return stageTime.running.expand.start.expand.transponder.differentiator
      }
    }
  }

  function formatMsToString (durationMs: number): string {
    let remainderMinutes = durationMs % 60000
    const minutes = ((durationMs - remainderMinutes) / 60000).toString().padStart(2, "0")
    
    let remainderSeconds = remainderMinutes % 1000
    const seconds = ((remainderMinutes - remainderSeconds) / 1000).toString().padStart(2, "0")

    const milliseconds = remainderSeconds.toPrecision(3).toString().padEnd(3, "0")

    return minutes + ":" + seconds + ":" + milliseconds
  }

  function getDurationStringFromStageTime (stageTime: RecordModel): string {
    const duration = getDurationMsFromStageTime(stageTime)
    return formatMsToString(duration)
  }

  function getDiffToNextString (): string | undefined {
    if (position == 1) {
      return "-"
    }

    if (stageTime.best == undefined) {
      return
    }

    const nextBestStageTime = stageTimes.values().toArray()[position-2].best

    if (nextBestStageTime == undefined) {
      return
    }

    const thisTime = getDurationMsFromStageTime(stageTime.best)
    const nextTime = getDurationMsFromStageTime(nextBestStageTime)

    return formatMsToString(thisTime - nextTime)
  }
</script>


<Table.Row>
  <Table.Cell>{position}</Table.Cell>
  <Table.Cell>
    {getDriverName()}
    {#if getDifferentiator()}
      <Badge variant="secondary">{getDifferentiator()}</Badge>
    {/if}
  </Table.Cell>
  {#if stageTime.best}
    <Table.Cell>{getDurationStringFromStageTime(stageTime.best)}</Table.Cell>
  {/if}
  {#if stageTime.last}
    <Table.Cell>{getDurationStringFromStageTime(stageTime.last)}</Table.Cell>
  {/if}
  <Table.Cell>{getDiffToNextString()}</Table.Cell>
  <Table.Cell>-<!-- TODO --></Table.Cell>
</Table.Row>
