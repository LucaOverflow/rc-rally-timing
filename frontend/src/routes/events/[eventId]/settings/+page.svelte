<script lang="ts">
  import EventSettings from '$lib/components/event-settings.svelte';
  import { globals } from '$lib/globals.svelte';
  import { Button } from "$lib/components/ui/button"
  import { pb } from '$lib/pb';
  import { toast } from 'svelte-sonner';

  let eventFormData = $state({
    name: globals.activeEvent?.name,
    description: globals.activeEvent?.description,
    start: globals.activeEvent?.start,
    end: globals.activeEvent?.end,
    location: globals.activeEvent?.location,
    start_order: globals.activeEvent?.start_order,
    start_mode: globals.activeEvent?.start_mode,
    signup_until: globals.activeEvent?.signup_until,
    max_drivers: globals.activeEvent?.max_drivers,
    allow_transponder_signup: globals.activeEvent?.allow_transponder_signup
    })

  function saveSettings () {
    if (globals.activeEvent == undefined) {
      toast.error("Couldn't save event settings. Please re-enter the event from the event selection.")
      return
    }

    pb.collection("events").update(globals.activeEvent.id, {
      ...eventFormData
    })
      .then((updatedEvent) => {
        toast.success("Event settings saved.")
        globals.activeEvent = updatedEvent
      })
      .catch(() => {
        toast.error("Couldn't save event settings.")
      })
  }
</script>

<div class="m-6">
  <EventSettings bind:formData={eventFormData} />
  <Button onclick={saveSettings} class="mt-4">Save</Button>
</div>