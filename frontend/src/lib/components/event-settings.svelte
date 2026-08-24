<script lang="ts">
  import * as Field from "$lib/components/ui/field"
  import { Separator } from "$lib/components/ui/separator"
  import { Input } from "$lib/components/ui/input"
  import { Textarea } from "$lib/components/ui/textarea"
  import DateTimePicker from './date-time-picker.svelte'
  import { Checkbox } from "$lib/components/ui/checkbox"
  import * as Select from "$lib/components/ui/select"

  let {
    formData = $bindable({
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
  } = $props()

  interface option {
    name: string,
    value: string,
    exclude?: string[]
  }

  const start_order_options: option[] = [
    {
      name: "Reactive",
      value: "reactive",
      exclude: [
        "fixed_time"
      ]
    },
    {
      name: "Ordered",
      value: "ordered"
    }
  ]

  const start_mode_options: option[] = [
    {
      name: "Passing",
      value: "passing"
    },
    {
      name: "Countdown",
      value: "countdown"
    },
    {
      name: "Fixed Time",
      value: "fixed_time"
    }
  ]

  let startModeSelectionDisabled = $derived(formData.start_order == "")

  function getOptionName (options: option[], value: string): string {
    const optionIndex = options.findIndex((option) => {
      return option.value == value
    })

    if (optionIndex == -1) {
      return ""
    }

    return options[optionIndex].name
  }

  function isOptionExcluded (optionsSource: option[], optionExcludeSource: string, optionExcludeTarget: string): boolean {
    const optionIndex = optionsSource.findIndex((option) => {
      return option.value == optionExcludeSource
    })

    if (optionIndex == -1) {
      return false
    }

    return optionsSource[optionIndex].exclude?.includes(optionExcludeTarget) || false
  }

  $effect(() => {
    if (isOptionExcluded(start_order_options, formData.start_order, formData.start_mode)) {
      formData.start_mode = ""
    }
  })
</script>

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
      <Input id="name" type="text" autocomplete="off" bind:value={formData.name} />
    </Field.Field>

    <Field.Field>
      <Field.Label for="description">Description</Field.Label>
      <Field.Description>
        Please provide all relevant info for participants. E.g. rules<br>
      </Field.Description>
      <Textarea id="description" bind:value={formData.description} />
    </Field.Field>

    <Field.Field>
      <Field.Label for="location">Location</Field.Label>
      <Textarea id="location" bind:value={formData.location} />
    </Field.Field>

    <Separator />

    <Field.Field>
      <Field.Label for="start">Start</Field.Label>
      <DateTimePicker id="start" bind:value={formData.start} />
    </Field.Field>

    <Field.Field>
      <Field.Label for="end">End</Field.Label>
      <DateTimePicker id="end" bind:value={formData.end} />
    </Field.Field>

    <Field.Field>
      <Field.Label for="signup_until">Signup until</Field.Label>
      <Field.Description>
        Leave this empty or after even end if you don't want to limit sign ups.
      </Field.Description>
      <DateTimePicker id="signup_until" bind:value={formData.signup_until} />
    </Field.Field>

    <Field.Field orientation="horizontal">
      <Checkbox id="transponder_signup" bind:checked={formData.allow_transponder_signup} />
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
      <Input id="max_drivers" type="number" bind:value={formData.max_drivers} />
    </Field.Field>

    <Separator />

    <Field.Field>
      <Field.Label for="start_order">Start Order</Field.Label>
      <Select.Root type="single" bind:value={formData.start_order}>
        <Select.Trigger id="start_order">{getOptionName(start_order_options, formData.start_order) || "Select a start order"}</Select.Trigger>
        <Select.Content>
          {#each start_order_options as option}
            <Select.Item value={option.value}>{option.name}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </Field.Field>

    <Field.Field>
      <Field.Label for="start_mode">Start Mode</Field.Label>
      <Select.Root type="single" bind:value={formData.start_mode} disabled={startModeSelectionDisabled}>
        <Select.Trigger id="start_mode">{getOptionName(start_mode_options, formData.start_mode) || "Select a start mode"}</Select.Trigger>
        <Select.Content>
          {#each start_mode_options as option}
            <Select.Item value={option.value} disabled={isOptionExcluded(start_order_options, formData.start_order, option.value)}>{option.name}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </Field.Field>

  </Field.Group>

</Field.Set>