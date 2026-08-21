<script lang="ts">
  import { getLocalTimeZone, Time, toCalendarDateTime } from "@internationalized/date";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import Calendar from "$lib/components/ui/calendar/calendar.svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import type { CalendarDate, CalendarDateTime } from "@internationalized/date";
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import { ChevronDown } from '@hugeicons/core-free-icons';
  
  let open = $state(false);
  let date = $state<CalendarDate | undefined>();
  let time = $state("10:00")

  let { 
    id: id = "",
    value: value = $bindable('')
   } = $props()

  $effect(() => {
    if (date == undefined) {
      return
    }

    let output = date.toDate(getLocalTimeZone())
    output.setHours(Number(time.split(":")[0]))
    output.setMinutes(Number(time.split(":")[1]))

    value = output.toISOString()
  })
</script>
 
<div class="flex gap-4">
 <div class="flex flex-col gap-3">
  <Popover.Root bind:open>
   <Popover.Trigger id={id}>
    {#snippet child({ props })}
     <Button
      {...props}
      variant="outline"
      class="w-32 justify-between font-normal"
     >
      {date
       ? date.toDate(getLocalTimeZone()).toLocaleDateString()
       : "Select date"}
      <HugeiconsIcon icon={ChevronDown} />
     </Button>
    {/snippet}
   </Popover.Trigger>
   <Popover.Content class="w-auto overflow-hidden p-0" align="start">
    <Calendar
     type="single"
     bind:value={date}
     onValueChange={() => {
      open = false;
     }}
     captionLayout="dropdown"
    />
   </Popover.Content>
  </Popover.Root>
 </div>
 <div class="flex flex-col gap-3">
  <Input
   type="time"
   id="{id}-time"
   step="3600"
   bind:value={time}
   class="appearance-none bg-background [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
  />
 </div>
</div>