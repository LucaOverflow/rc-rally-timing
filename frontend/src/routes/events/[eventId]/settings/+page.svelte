<script lang="ts">
  import EventSettings from '$lib/components/event-settings.svelte';
  import { globals } from '$lib/globals.svelte';

  // TODO Fix date parsing

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
</script>

<div class="m-6">
  <EventSettings bind:formData={eventFormData} />
</div>