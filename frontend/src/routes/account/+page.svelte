<script lang="ts">
  import * as Alert from "$lib/components/ui/alert/index"
  import { Input } from '$lib/components/ui/input/index'
  import { Label } from '$lib/components/ui/label/index'
  import { Button } from "$lib/components/ui/button/index"
  import { pb } from '$lib/pb'
  import { Alert01Icon, InformationCircleIcon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon } from '@hugeicons/svelte';
  import { toast } from 'svelte-sonner';
  
  let initialPersonalDataFormData = {
    name: pb.authStore.record?.name,
    email: pb.authStore.record?.email,
    nameErrorMessage: '',
    emailErrorMessage: '',
    nameChanged: false,
    emailChangeRequested: false
  }
  let personalDataFormData = $state(initialPersonalDataFormData)
  let enablePersonalDataFormSaveButton = $derived(personalDataFormData.name != initialPersonalDataFormData.name || personalDataFormData.email != initialPersonalDataFormData.email)

  let changePasswordFormData = $state({
    oldPassword: '',
    password: '',
    passwordConfirm: '',
    errorMessage: '',
    passwordChanged: false
  })
  let enableChangePasswordFormSaveButton = $derived(changePasswordFormData.oldPassword != '' && changePasswordFormData.password != '' && changePasswordFormData.passwordConfirm != '')

  function copyUserIdToClipboard () {
    if (pb.authStore.record?.id) {
      navigator.clipboard.writeText(pb.authStore.record.id)
        .then(() => {
          toast.success('Copied User ID to clipboard')
        })
        .catch(() => {
          toast.error("Couldn't copy User ID to clipboard")
        })
    } else {
      toast.error('Please Login first')
    }
  }

  function savePersonalData (event: Event) {
    event.preventDefault()

    personalDataFormData.nameErrorMessage = ''
    personalDataFormData.nameChanged = false
    personalDataFormData.emailErrorMessage = ''
    personalDataFormData.emailChangeRequested = false

    if (personalDataFormData.name != initialPersonalDataFormData.name) {
      const userId = pb.authStore.record?.id || ''
      pb.collection('users').update(userId, {
        name: personalDataFormData.name
      })
        .then(() => {
          personalDataFormData.nameChanged = true
          initialPersonalDataFormData.name = personalDataFormData.name
        })
        .catch((error: Error) => {
          personalDataFormData.nameErrorMessage = error.message
        })
    }

    if (personalDataFormData.email != initialPersonalDataFormData.email) {
      pb.collection('users').requestEmailChange(personalDataFormData.email)
        .then(() => {
          personalDataFormData.emailChangeRequested = true
          initialPersonalDataFormData.email = personalDataFormData.email
        })
        .catch((error: Error) => {
          personalDataFormData.emailErrorMessage = error.message
        })
    }
  }

  function changePassword (event: Event) {
    event.preventDefault()

    changePasswordFormData.errorMessage = ''
    changePasswordFormData.passwordChanged = false

    const userId = pb.authStore.record?.id || ''
    const userMail = pb.authStore.record?.email
    pb.collection('users').update(userId, {
      oldPassword: changePasswordFormData.oldPassword,
      password: changePasswordFormData.password,
      passwordConfirm: changePasswordFormData.passwordConfirm
    })
      .then(() => {
        pb.collection('users').authWithPassword(userMail, changePasswordFormData.password)
        changePasswordFormData.oldPassword = ''
        changePasswordFormData.password = ''
        changePasswordFormData.passwordConfirm = ''
        changePasswordFormData.passwordChanged = true
      })
      .catch((error: Error) => {
        changePasswordFormData.errorMessage = error.message
      })
  }
</script>

<div class="flex ml-10 mr-10">
  <div class="flex-2">
    <h1 class="pb-6">Personal Data</h1>

    <Button onclick={copyUserIdToClipboard} class="mb-6">Copy User ID to Clipboard</Button>

    <form onsubmit={savePersonalData} class="space-y-1">
      <Label for="name">Name</Label>
      <Input id="name" name="name" type="text" bind:value={personalDataFormData.name} class="mb-6" />
      {#if personalDataFormData.nameErrorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to change Name</Alert.Title>
          <Alert.Description>
            {personalDataFormData.nameErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

      <Label for="email">E-Mail</Label>
      <Input id="email" name="email" type="email" bind:value={personalDataFormData.email} class="mb-6" />
      {#if personalDataFormData.emailErrorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to change E-Mail</Alert.Title>
          <Alert.Description>
            {personalDataFormData.emailErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

      <Button type="submit" disabled={!enablePersonalDataFormSaveButton}>Save</Button>

      {#if personalDataFormData.nameChanged}
        <Alert.Root>
            <HugeiconsIcon icon={InformationCircleIcon} />
            <Alert.Title>New Name saved</Alert.Title>
          </Alert.Root>
      {/if}

      {#if personalDataFormData.emailChangeRequested}
        <Alert.Root>
            <HugeiconsIcon icon={InformationCircleIcon} />
            <Alert.Title>E-Mail change requested</Alert.Title>
            <Alert.Description>
              Please look into your inbox to verify your new E-Mail.
            </Alert.Description>
          </Alert.Root>
      {/if}
    </form>
  </div>
  <div class="flex-1"></div>
  <div class="flex-2">
    <h1 class="pb-6">Change Password</h1>

    <form onsubmit={changePassword} class="space-y-1">
      <Label for="oldPassword">Old password</Label>
      <Input id="oldPassword" name="oldPassword" type="password" bind:value={changePasswordFormData.oldPassword} class="mb-6" />

      <Label for="password">New password</Label>
      <Input id="password" name="password" type="password" bind:value={changePasswordFormData.password} class="mb-6" />

      <Label for="passwordConfirm">Confirm new password</Label>
      <Input id="passwordConfirm" name="passwordConfirm" type="password" bind:value={changePasswordFormData.passwordConfirm} class="mb-6" />

      <Button type="submit" disabled={!enableChangePasswordFormSaveButton}>Save</Button>

      {#if changePasswordFormData.passwordChanged}
        <Alert.Root>
            <HugeiconsIcon icon={InformationCircleIcon} />
            <Alert.Title>Your passwort has been changed</Alert.Title>
          </Alert.Root>
      {/if}

      {#if changePasswordFormData.errorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to change Password</Alert.Title>
          <Alert.Description>
            {changePasswordFormData.errorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}
    </form>
  </div>
</div>