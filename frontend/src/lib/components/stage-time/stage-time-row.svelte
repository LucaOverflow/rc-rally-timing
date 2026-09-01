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

  export function getDurationStringFromStageTime (stageTime: RecordModel): string {
    const duration = getDurationMsFromStageTime(stageTime)

    let remainderMinutes = duration % 60000
    const minutes = ((duration - remainderMinutes) / 60000).toString().padStart(2, "0")
    
    let remainderSeconds = remainderMinutes % 1000
    const seconds = ((remainderMinutes - remainderSeconds) / 1000).toString().padStart(2, "0")

    const milliseconds = remainderSeconds.toPrecision(3).toString().padEnd(3, "0")

    return minutes + ":" + seconds + ":" + milliseconds
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
  <Table.Cell>-<!-- TODO --></Table.Cell>
  <Table.Cell>-<!-- TODO --></Table.Cell>
</Table.Row>
