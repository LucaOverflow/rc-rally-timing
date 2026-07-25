<script module>
  export interface stageTime {
    best: RecordModel | undefined,
    last: RecordModel | undefined,
    running: RecordModel | undefined
  }
</script>

<script lang="ts">
  import type { RecordModel } from 'pocketbase'

  let { stageTime: stageTime } = $props()
</script>

<p>{stageTime[0]}</p>
