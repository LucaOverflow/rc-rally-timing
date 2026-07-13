<script lang="ts">
  import { Badge } from "$lib/components/ui/badge"
  import { Button } from "$lib/components/ui/button"
  import { ScrollArea } from "$lib/components/ui/scroll-area"
  import { Input } from "$lib/components/ui/input"
  import { Textarea } from "$lib/components/ui/textarea"
  import { Label } from "$lib/components/ui/label"
  import { Spinner } from "$lib/components/ui/spinner"
  import * as Item from "$lib/components/ui/item"
  import * as Dialog from "$lib/components/ui/dialog"
  import * as Alert from "$lib/components/ui/alert"
  import * as InputGroup from "$lib/components/ui/input-group"
  import { Alert01Icon, Cancel01Icon, PlusSignIcon } from '@hugeicons/core-free-icons'
  import { HugeiconsIcon } from '@hugeicons/svelte'
  import { onMount } from 'svelte'
  import { pb } from '$lib/pb'
  import type { RecordModel } from 'pocketbase'
  import { toast } from 'svelte-sonner'

  let decoder: RecordModel[] = $state([])
  let initialDecoderState: RecordModel[]
  let decoderStatusColor: string[] = $state([])
  let saveTimeout: number

  let openAddDecoderPopup = $state(false)
  let addDecoderFormData = $state({
    name: undefined,
    personal_notes: ''
  })
  let addDecoderErrorMessage = $state('')

  let addDecoderOwnerId: string[] = $state([])

  onMount(() => {
    requestDecoder()
  })

  function requestDecoder () {
    pb.collection('decoder').getFullList({
      filter: pb.filter('owner ~ {:owner}', {owner: pb.authStore.record?.id}),
      expand: 'owner'
    })
      .then((data: RecordModel[]) => {
        decoder = data
        initialDecoderState = data
        refreshDecoderStatus()
      })
  }

  async function refreshDecoderStatus () {
    for (let i = 0; i < decoder.length; ++i) {
      const decoderPassings = await pb.collection('passings').getList(1, 1, {
        filter: pb.filter('decoder = {:decoder}', {decoder: decoder[i].id})
      })

      if (decoderPassings.items.length > 0) {
        decoderStatusColor[i] = 'green'
      } else {
        decoderStatusColor[i] = 'gray'
      }
    }
  }

  function addDecoder () {
    addDecoderErrorMessage = ''

    pb.collection('decoder').create({
      ...addDecoderFormData,
      owner: [pb.authStore.record?.id]
    })
      .then(() => {
        openAddDecoderPopup = false
        requestDecoder()
      })
      .catch((error: Error) => {
        addDecoderErrorMessage = error.message
      })
  }

  function queueNoteSave () {
    clearTimeout(saveTimeout)
    saveTimeout = setTimeout(saveNotes, 1000)
  }

  function saveNotes () {
    let wereAllNotesSaved = true

    for (let i = 0; i < decoder.length; ++i) {
      if (decoder[i].personal_notes != initialDecoderState[i].personal_notes) {
        pb.collection('decoder').update(decoder[i].id, {
          personal_notes: decoder[i].personal_notes
        })
          .then(() => {
            initialDecoderState[i] = {...decoder[i]}
          })
          .catch((_error: Error) => {
            wereAllNotesSaved = false
          })
          .finally(() => {
            if (wereAllNotesSaved) {
              toast.success("Notes saved")
            } else {
              toast.error("Notes couldn't all be saved")
            }
          })
      }
    }
  }

  function removeOwnerFromDecoder (decoder_id: string, user_id: string) {
    pb.collection('decoder').update(decoder_id, {
      'owner-': user_id
    })
      .then(() => {
        toast.success('Removed owner from Decoder')
        requestDecoder()
      })
      .catch((_error: Error) => {
        toast.error("couldn't remove owner from Decoder")
      })
  }

  async function addDecoderOwner (event: Event, decoder_index: number) {
    event.preventDefault()

    pb.collection('decoder').update(decoder[decoder_index].id, {
      'owner+': addDecoderOwnerId[decoder_index]
    })
      .then(() => {
        toast.success('Added owner to Decoder')
        addDecoderOwnerId[decoder_index] = ''
        requestDecoder()
      })
      .catch(() => {
        toast.error("Couldn't add owner to Decoder")
      })
  }

  function copyDecoderIdToClipboard (decoder_id: string) {
    navigator.clipboard.writeText(decoder_id)
      .then(() => {
        toast.success('Copied Decoder access key to clipboard')
      })
      .catch(() => {
        toast.error("Couldn't copy Decoder access key to clipboard")
      })
  }

  function deleteDecoder (recordId: string) {
    pb.collection('decoder').delete(recordId)
      .then(() => {
        toast.success("Decoder removed")
        requestDecoder()
      })
      .catch(() => {
        toast.error("Couldn't remove Decoder")
      })
  }
</script>

<Button onclick={() => {openAddDecoderPopup = true}} class="mt-2">
  <HugeiconsIcon icon={PlusSignIcon} />
  Add
</Button>

<ScrollArea class="m-6">
  {#each decoder as thisDecoder, i}
    <Item.Root variant="outline" class="mb-2">
      <Item.Content class="flex-row space-x-1 items-center">
        <div class="w-1/3 mb-auto">
          <Item.Title>
            {#if decoderStatusColor[i]}
              <div class="w-3 h-3 rounded-full" style="background-color: {decoderStatusColor[i]}"></div>
            {:else}
              <Spinner />
            {/if}
            {thisDecoder.name}
            {#if thisDecoder.type != ''}
              <Badge variant="secondary">{thisDecoder.type}</Badge>
            {/if}
          </Item.Title>
          <Textarea placeholder="Personal Note" bind:value={thisDecoder.personal_notes} oninput={queueNoteSave} class="mt-2" />
        </div>
        <div class="flex flex-col gap-2 ml-auto">
          <h1>Owner</h1>
          {#each decoder[i].expand?.owner as owner}
            <Badge variant="outline" class="p-0">
              <span class="ml-2">{owner.name}</span>
              <Button variant="destructive" onclick={() => { removeOwnerFromDecoder(thisDecoder.id, owner.id) }} class="w-5 h-5 m-0 rounded-full border-none">
                <HugeiconsIcon icon={Cancel01Icon} />
              </Button>
            </Badge>
          {/each}
          <form onsubmit={(event: Event) => { addDecoderOwner(event, i) }}>
            <InputGroup.Root>
              <InputGroup.Input type="text" bind:value={addDecoderOwnerId[i]} placeholder='User ID' />
              <InputGroup.Button variant="default" type="submit" class="mr-1">Add</InputGroup.Button>
            </InputGroup.Root>
          </form>
        </div>
        <div class="ml-auto flex flex-col gap-2">
          <Button onclick={() => { copyDecoderIdToClipboard(thisDecoder.id) }}>Copy access key</Button>
          <Button variant="secondary">Show logs</Button>
          <Button variant="destructive" onclick={() => {deleteDecoder(thisDecoder.id)}}>Remove</Button>
        </div>
      </Item.Content>
    </Item.Root>
  {/each}
</ScrollArea>

<!-- Add Decoder Popup -->
<Dialog.Root bind:open={openAddDecoderPopup}>
<Dialog.Content>
  <form onsubmit={addDecoder}>
    <Dialog.Header>

      <Label for="name">Name</Label>
      <Input id="name" name="name" type="text" bind:value={addDecoderFormData.name} />

      <Label for="personal_notes">Personal Notes</Label>
      <Input id="personal_notes" name="personal_notes" type="text" bind:value={addDecoderFormData.personal_notes} />

      {#if addDecoderErrorMessage != ''}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to add Decoder</Alert.Title>
          <Alert.Description>
            {addDecoderErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button type="submit">Add</Button>
    </Dialog.Footer>
  </form>
</Dialog.Content>
</Dialog.Root>