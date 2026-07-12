<script lang="ts">
  import { Badge } from "$lib/components/ui/badge"
  import { Button } from "$lib/components/ui/button"
  import { ScrollArea } from "$lib/components/ui/scroll-area"
  import { Input } from "$lib/components/ui/input"
  import { Label } from "$lib/components/ui/label"
  import { Spinner } from "$lib/components/ui/spinner"
  import * as Item from "$lib/components/ui/item"
  import * as Dialog from "$lib/components/ui/dialog"
  import * as Alert from "$lib/components/ui/alert"
  import { Alert01Icon, PlusSignIcon } from '@hugeicons/core-free-icons'
  import { HugeiconsIcon } from '@hugeicons/svelte'
  import { onMount } from 'svelte'
  import { pb } from '$lib/pb'
  import type { RecordModel } from 'pocketbase'
  import { toast } from 'svelte-sonner'

  let transponder: RecordModel[] = $state([])
  let initialTransponderState: RecordModel[]
  let transponderStatusColor: string[] = $state([])
  let saveTimeout: number

  let openAddTransponderPopup = $state(false)
  let addTransponderFormData = $state({
    id: undefined,
    personal_notes: ''
  })
  let addTransponderErrorMessage = $state('')

  onMount(() => {
    requestTransponder()
  })

  function requestTransponder () {
    pb.collection('transponder').getFullList({
      filter: pb.filter('owner = {:owner}', {owner: pb.authStore.record?.id})
    })
      .then((data: RecordModel[]) => {
        transponder = data
        initialTransponderState = data
        refreshTransponderStatus()
      })
  }

  async function refreshTransponderStatus () {
    for (let i = 0; i < transponder.length; ++i) {
      const transponderPassings = await pb.collection('passings').getList(1, 1, {
        filter: pb.filter('transponder = {:transponder}', {transponder: transponder[i].id})
      })

      if (transponderPassings.items.length > 0) {
        transponderStatusColor[i] = 'green'
      } else {
        transponderStatusColor[i] = 'gray'
      }
    }
  }

  function addTransponder () {
    addTransponderErrorMessage = ''

    pb.collection('transponder').create({
      ...addTransponderFormData,
      owner: pb.authStore.record?.id
    })
      .then(() => {
        openAddTransponderPopup = false
        requestTransponder()
      })
      .catch((error: Error) => {
        addTransponderErrorMessage = error.message
      })
  }

  function queueNoteSave () {
    clearTimeout(saveTimeout)
    saveTimeout = setTimeout(saveNotes, 1000)
  }

  function saveNotes () {
    let wereAllNotesSaved = true

    for (let i = 0; i < transponder.length; ++i) {
      if (transponder[i].personal_notes != initialTransponderState[i].personal_notes) {
        pb.collection('transponder').update(transponder[i].id, transponder[i])
          .then(() => {
            initialTransponderState[i] = {...transponder[i]}
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

  function deleteTransponder (recordId: string) {
    pb.collection('transponder').delete(recordId)
      .then(() => {
        toast.success("Transponder removed")
        requestTransponder()
      })
      .catch(() => {
        toast.error("Couldn't remove Transponder")
      })
  }
</script>

<Button onclick={() => {openAddTransponderPopup = true}} class="mt-2">
  <HugeiconsIcon icon={PlusSignIcon} />
  Add
</Button>

<ScrollArea class="m-6">
  {#each transponder as thisTransponder, i}
    <Item.Root variant="outline" class="mb-2">
      <Item.Media>
        {#if transponderStatusColor[i]}
          <div class="w-3 h-3 rounded-full" style="background-color: {transponderStatusColor[i]}"></div>
        {:else}
          <Spinner />
        {/if}
      </Item.Media>
      <Item.Content class="flex-row space-x-1 items-center">
        <Item.Title>{thisTransponder.id}</Item.Title>
        {#if thisTransponder.type != ''}
          <Badge variant="secondary">{thisTransponder.type}</Badge>
        {/if}
        <Input type="text" placeholder="Personal Note" bind:value={thisTransponder.personal_notes} oninput={queueNoteSave} class="ml-auto w-1/2" />
        <Button variant="destructive" class="ml-auto" onclick={() => {deleteTransponder(thisTransponder.id)}}>Remove</Button>
      </Item.Content>
    </Item.Root>
  {/each}
</ScrollArea>

<!-- Add Transponder Popup -->
<Dialog.Root bind:open={openAddTransponderPopup}>
<Dialog.Content>
  <form onsubmit={addTransponder}>
    <Dialog.Header>

      <Label for="id">Transponder ID</Label>
      <Input id="id" name="id" type="number" bind:value={addTransponderFormData.id} />

      <Label for="personal_notes">Personal Notes</Label>
      <Input id="personal_notes" name="personal_notes" type="text" bind:value={addTransponderFormData.personal_notes} />

      {#if addTransponderErrorMessage != ''}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to add Transponder</Alert.Title>
          <Alert.Description>
            {addTransponderErrorMessage}
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