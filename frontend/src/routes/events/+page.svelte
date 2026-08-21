<script lang="ts">
  import { goto } from '$app/navigation'
  import { Button } from "$lib/components/ui/button"
  import * as Item from '$lib/components/ui/item'
  import * as Dialog from "$lib/components/ui/dialog"
  import * as Alert from "$lib/components/ui/alert"
  import * as Field from "$lib/components/ui/field"
  import * as Select from "$lib/components/ui/select"
  import { Checkbox } from "$lib/components/ui/checkbox"
  import { Textarea } from "$lib/components/ui/textarea"
  import { Separator } from "$lib/components/ui/separator"
  import { ScrollArea } from '$lib/components/ui/scroll-area'
  import { Input } from "$lib/components/ui/input"
  import { globals } from '$lib/globals.svelte'
  import { pb } from '$lib/pb'
  import { Alert01Icon, PlusSignIcon, StopWatchIcon } from '@hugeicons/core-free-icons'
  import { HugeiconsIcon } from '@hugeicons/svelte'
  import type { ListResult, RecordModel } from 'pocketbase'
  import { onMount } from 'svelte'
  import { toast } from 'svelte-sonner'
  import DateTimePicker from '$lib/components/date-time-picker.svelte';

  let events: RecordModel[] = $state([])

  let openAddEventPopup = $state(false)
  let addEventFormData = $state({
    name: "",
    description: "",
    start: "",
    end: "",
    location: "",
    start_order: "",
    start_mode: "",
    signup_until: "",
    max_drivers: 50,
    allow_transponder_signup: true
  })
  let addEventErrorMessage = $state("")

  onMount(() => {
    requestEvents()
  })

  function requestEvents () {
    pb.collection('events').getList(1, 25, {
      filter: pb.filter('organizer ~ {:user}', {user: pb.authStore.record?.id}),
      sort: '-created'
    })
      .then((result: ListResult<RecordModel>) => {
        events = result.items
      })
      .catch(() => {
        toast.error("Couldn't load Events")
      })
  }

  function addEvent () {
  }

  function quickAddPractice () {
    const date = new Date()

    pb.collection('events').create({
      name: "Quick practice " + date.toLocaleDateString(),
      start_order: "reactive",
      start_mode: "passing",
      organizer: pb.authStore.record?.id,
      allow_transponder_signup: true
    })
      .then(() => {
        requestEvents()
      })
      .catch(() => {
        toast.error("Couldn't add event")
      })
  }

  function deleteEvent (event: string) {
    pb.collection('events').delete(event)
      .then(() => {
        requestEvents()
        if (globals.activeEvent?.id == event) {
          globals.activeEvent = undefined
        }
      })
      .catch(() => {
        toast.error("Couldn't remove event")
      })
  }
</script>

<Button onclick={() => {openAddEventPopup = true}} class="mt-2">
  <HugeiconsIcon icon={PlusSignIcon} />
  Add
</Button>

<Button onclick={quickAddPractice} class="mt-2">
  <HugeiconsIcon icon={StopWatchIcon} />
  Quick Add Practice
</Button>

<ScrollArea class="m-6">
  {#each events as event, i}
    <Item.Root variant="outline" class="mb-2">
      <Item.Content class="flex-row space-x-1 items-center">
        <Item.Title>{event.name}</Item.Title>
        <div class="flex flex-col gap-1 ml-auto">
          <Button onclick={() => {goto('/events/' + event.id + '/schedule')}}>Select</Button>
          <Button variant="destructive" onclick={() => {deleteEvent(event.id)}}>Remove</Button>
        </div>
      </Item.Content>
    </Item.Root>
  {/each}
</ScrollArea>

<!-- TODO Add pagination navigation here -->

<!-- Add Event Popup -->
<Dialog.Root bind:open={openAddEventPopup}>
  <Dialog.Content class="flex max-h-[90%] flex-col gap-0">
    <Dialog.Header class="overflow-y-auto">
      <Field.Set>

        <Field.Legend>
          Add a new Event
        </Field.Legend>
        <Field.Description>
          You can still edit all these settings after creation.
        </Field.Description>

        <Separator />

        <Field.Group>

          <Field.Field>
            <Field.Label for="name">Name</Field.Label>
            <Input id="name" type="text" autocomplete="off" bind:value={addEventFormData.name} />
          </Field.Field>

          <Field.Field>
            <Field.Label for="description">Description</Field.Label>
            <Field.Description>
              Please provide all relevant info for participants. E.g. rules<br>
            </Field.Description>
            <Textarea id="description" bind:value={addEventFormData.description} />
          </Field.Field>

          <Field.Field>
            <Field.Label for="location">Location</Field.Label>
            <Textarea id="location" bind:value={addEventFormData.location} />
          </Field.Field>

          <Separator />

          <Field.Field>
            <Field.Label for="start">Start</Field.Label>
            <DateTimePicker id="start" bind:value={addEventFormData.start} />
          </Field.Field>

          <Field.Field>
            <Field.Label for="end">End</Field.Label>
            <DateTimePicker id="end" bind:value={addEventFormData.end} />
          </Field.Field>

          <Field.Field>
            <Field.Label for="signup_until">Signup until</Field.Label>
            <Field.Description>
              Leave this empty or after even end if you don't want to limit sign ups.
            </Field.Description>
            <DateTimePicker id="signup_until" bind:value={addEventFormData.signup_until} />
          </Field.Field>

          <Field.Field orientation="horizontal">
            <Checkbox id="transponder_signup" bind:checked={addEventFormData.allow_transponder_signup} />
            <Field.Label for="transponder_signup">Allow Transponder Signup</Field.Label>
            <Field.Description>
              If this is enabled drivers will automatically be signed up if they pass through a decoder connected to the event.
            </Field.Description>
          </Field.Field>

          <Field.Field>
            <Field.Label for="max_drivers">Max. Drivers</Field.Label>
            <Field.Description>
              This is the overall event size limit. Additionally you'll be able to limit the class sizes induvitually after event creation.
            </Field.Description>
            <Input id="max_drivers" type="number" bind:value={addEventFormData.max_drivers} />
          </Field.Field>

          <Separator />

          <Field.Field>
            <Field.Label for="start_order">Start Order</Field.Label>
            <Select.Root type="single" bind:value={addEventFormData.start_order}>
              <Select.Trigger id="start_order">{addEventFormData.start_order || "Select a start order"}</Select.Trigger>
              <Select.Content>
                <Select.Item value="reactive">Reactive</Select.Item>
                <Select.Item value="ordered">Ordered</Select.Item>
              </Select.Content>
            </Select.Root>
          </Field.Field>

          <Field.Field>
            <Field.Label for="start_mode">Start Mode</Field.Label>
            <Select.Root type="single" bind:value={addEventFormData.start_mode}>
              <Select.Trigger id="start_mode">{addEventFormData.start_mode || "Select a start mode"}</Select.Trigger>
              <Select.Content>
                <Select.Item value="passing">Passing</Select.Item>
                <Select.Item value="countdown">Countdown</Select.Item>
                <Select.Item value="fixed_time">Fixed Time</Select.Item>
              </Select.Content>
            </Select.Root>
          </Field.Field>

        </Field.Group>

      </Field.Set>

      {#if addEventErrorMessage != ''}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to add Event</Alert.Title>
          <Alert.Description>
            {addEventErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button onclick={addEvent}>Add</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>