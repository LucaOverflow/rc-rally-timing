<script lang="ts">
  import * as Select from "$lib/components/ui/select"
  import { pb } from '$lib/pb';
  import type { ClientResponseError, RecordModel } from 'pocketbase';
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';

  let { value: selectedDecoder = $bindable('') } = $props()
  let selectedDecoderName = $derived(getDecoderName(selectedDecoder))

  interface inUseDecoder {
    decoder: RecordModel,
    usedBy: RecordModel
  }

  let decoder = $state({
    available: [] as RecordModel[],
    inUse: [] as inUseDecoder[]
  })

  onMount(() => {
    requestDecoder()
  })

  function requestDecoder () {
    decoder.available = []
    decoder.inUse = []

    pb.collection('decoder').getFullList({
      filter: pb.filter('owner ~ {:currentUser}', {currentUser: pb.authStore.record?.id}),
      sort: '+name',
      requestKey: null
    })
      .then((decoderResult) => {
        for (const thisDecoder of decoderResult) {
          pb.collection('stages').getFirstListItem(
            pb.filter('start_decoder = {:decoder} || stop_decoder = {:decoder}', {decoder: thisDecoder.id}),
            {
              requestKey: null
            }
          )
            .then((stagesResult) => {
              decoder.inUse.push({
                decoder: thisDecoder,
                usedBy: stagesResult
              })
            })
            .catch((error: ClientResponseError) => {
              if (error.status == 404) {
                decoder.available.push(thisDecoder)
              } else {
                toast.error("Couldn't check if Decoder " + thisDecoder.name + " is already in use")
              }
            })
        }
      })
      .catch((error: ClientResponseError) => {
        if (error.status == 404) {
          toast.error("You don't own any decoder")
        } else {
          toast.error("Couldn't request decoder")
        }
      })
  }

  function getDecoderName (decoderId: string): string {
    const foundAvailable = decoder.available.find((element) => {
      return element.id == decoderId
    })

    if (foundAvailable != undefined) {
      return foundAvailable.name
    }

    const foundInUse = decoder.inUse.find((element) => {
      return element.decoder.id == decoderId
    })

    if (foundInUse != undefined) {
      return foundInUse.decoder.name
    }

    return ''
  }
</script>

<Select.Root type="single" bind:value={selectedDecoder}>
  <Select.Trigger>
    {#if selectedDecoder == ''}
      Select a Decoder
    {:else}
      {selectedDecoderName}
    {/if}
  </Select.Trigger>
  <Select.Content>
    <Select.Group>
      {#if decoder.available.length > 0}
        <Select.Label>Available Decoder</Select.Label>
      {:else}
        <Select.Label>No available Decoder</Select.Label>
      {/if}
      {#each decoder.available as thisDecoder}
        <Select.Item value={thisDecoder.id}>{thisDecoder.name}</Select.Item>
      {/each}
    </Select.Group>
    {#if decoder.inUse.length > 0}
      <Select.Group>
        <Select.Label>Used Decoder</Select.Label>
        {#each decoder.inUse as thisDecoder}
          <Select.Item value={thisDecoder.decoder.id} disabled>{thisDecoder.decoder.name} - {thisDecoder.usedBy.name}</Select.Item>
        {/each}
      </Select.Group>
    {/if}
  </Select.Content>
</Select.Root>