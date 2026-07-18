<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index"
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index"
  import * as Dialog from "$lib/components/ui/dialog/index"
  import * as Alert from "$lib/components/ui/alert/index"
  import { Button } from "$lib/components/ui/button/index"
  import { Label } from "$lib/components/ui/label/index"
  import { Input } from "$lib/components/ui/input/index"
  import {
    Alert01Icon,
    Calendar03Icon,
    CalendarClockIcon,
    ChevronUpIcon,
    InformationCircleIcon,
    InternetAntenna02Icon,
    RankingIcon,
    StopWatchIcon,
    Timer02Icon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon, type IconSvgElement } from '@hugeicons/svelte';
  import { pb } from '$lib/pb';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';

  interface MenuItem {
    title: string,
    url: string,
    icon: IconSvgElement
  }

  const driverItems: MenuItem[] = [
    {
      title: "Events",
      url: "events",
      icon: Calendar03Icon
    },
    {
      title: "Transponder",
      url: "transponder",
      icon: InternetAntenna02Icon
    },
    {
      title: "Decoder",
      url: "decoder",
      icon: StopWatchIcon
    }
  ]

  const eventItems: MenuItem[] = [
    {
      title: "Schedule",
      url: "#",
      icon: CalendarClockIcon
    },
    {
      title: "Leaderboards",
      url: "#",
      icon: RankingIcon
    },
    {
      title: "Penalties Tracker",
      url: "#",
      icon: Timer02Icon
    }
  ]

  const isEventSelected = false // TODO Make dynamic

  let isLoggedIn = $state(pb.authStore.isValid)
  let openLoginPopup = $state(false)
  let openRegisterPopup = $state(false)
  let openResetPasswordPopup = $state(false)
  let resetPasswordRequested = $state(false)

  let formData = $state({
    name: '',
    email: '',
    password: '',
    passwordConfirm: '',
    loginErrorMessage: '',
    registerErrorMessage: '',
    resetPasswordErrorMessage: ''
  })

  onMount(() => {
    pb.authStore.onChange(() => {
      isLoggedIn = pb.authStore.isValid
    })
  })

  function resetFormData () {
    formData.name = ''
    formData.email = ''
    formData.password = ''
    formData.passwordConfirm = ''
    formData.loginErrorMessage = ''
    formData.registerErrorMessage = ''
    formData.resetPasswordErrorMessage = ''
  }

  function login (event: Event) {
    event.preventDefault()

    formData.loginErrorMessage = ''
    
    pb.collection('users').authWithPassword(formData.email, formData.password)
      .then(() => {
        openLoginPopup = false
        resetFormData()
      })
      .catch((error: Error) => {
        formData.loginErrorMessage = error.message
      })
  }
  
  function register (event: Event) {
    event.preventDefault()

    formData.registerErrorMessage = ''

    pb.collection('users').create({
      ...formData,
      emailVisibility: true
    })
      .then(() => {
        pb.collection('users').requestVerification(formData.email)
        login(event)
        openRegisterPopup = false
        resetFormData()
      })
      .catch((error: Error) => {
        formData.registerErrorMessage = error.message
      })
  }

  function resetPassword (event: Event) {
    event.preventDefault()

    formData.resetPasswordErrorMessage = ''

    pb.collection('users').requestPasswordReset(formData.email)
      .then(() => {
        resetPasswordRequested = true
        resetFormData()
      })
      .catch((error: Error) => {
        formData.resetPasswordErrorMessage = error.message
      })
  }

  function signout () {
    pb.authStore.clear()
  }
</script>

<Sidebar.Root>

  <Sidebar.Content>
    {#if isLoggedIn}
      <Sidebar.Group>
        <Sidebar.GroupLabel>Personal</Sidebar.GroupLabel>
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each driverItems as item (item.title)}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton isActive={page.url.pathname.endsWith(item.url)}>
                  {#snippet child({ props })}
                    <a href={item.url} {...props}>
                        <HugeiconsIcon icon={item.icon} />
                        <span>{item.title}</span>
                      </a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/each}
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    {/if}

    {#if isEventSelected}
      <Sidebar.Group>
        <Sidebar.GroupLabel>Event</Sidebar.GroupLabel> <!-- TODO Add Event Name. Always show currently running Event as default. -->
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each eventItems as item (item.title)}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton isActive={page.url.pathname.endsWith(item.url)}>
                  {#snippet child({ props })}
                    <a href={item.url} {...props}>
                        <HugeiconsIcon icon={item.icon} />
                        <span>{item.title}</span>
                      </a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/each}
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    {/if}
  </Sidebar.Content>

  <Sidebar.Footer>
    <Sidebar.Menu>
        <Sidebar.MenuItem>
          {#if isLoggedIn}
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                {#snippet child({ props })}
                  <Sidebar.MenuButton
                    {...props}
                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                  >
                    {pb.authStore.record?.name}
                    <HugeiconsIcon icon={ChevronUpIcon} class="ms-auto" />
                  </Sidebar.MenuButton>
                {/snippet}
              </DropdownMenu.Trigger>
              <DropdownMenu.Content
                side="top"
                class="w-(--bits-dropdown-menu-anchor-width)"
              >
                <DropdownMenu.Item onclick={() => { goto('account') }}>
                  <span>Account</span>
                </DropdownMenu.Item>
                <DropdownMenu.Item onclick={signout}>
                  <span>Sign out</span>
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          {:else}
            <Sidebar.Menu>
              <Sidebar.MenuItem>
                <div class="flex items-center gap-2">
                  <Sidebar.MenuButton class="justify-center">
                    {#snippet child({ props })}
                      <a onclick={() => { openLoginPopup = true }} {...props}>
                        <span>Log in</span>
                      </a>
                    {/snippet}
                  </Sidebar.MenuButton>
                  <Sidebar.MenuButton class="justify-center">
                    {#snippet child({ props })}
                      <a onclick={() => { openRegisterPopup = true }} {...props}>
                        <span>Register</span>
                      </a>
                    {/snippet}
                  </Sidebar.MenuButton>
                </div>
              </Sidebar.MenuItem>
            </Sidebar.Menu>
          {/if}
        </Sidebar.MenuItem>
      </Sidebar.Menu>
  </Sidebar.Footer>

</Sidebar.Root>

<!-- Login Popup -->
<Dialog.Root bind:open={openLoginPopup}>
<Dialog.Content>
  <form onsubmit={login}>
    <Dialog.Header>

      <Label for="email">E-Mail</Label>
      <Input id="email" name="email" type="email" bind:value={formData.email} />

      <Label for="password">Password</Label>
      <Input id="password" name="password" type="password" bind:value={formData.password} />

      {#if formData.loginErrorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to Login</Alert.Title>
          <Alert.Description>
            {formData.loginErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="link" class="mr-auto text-muted-foreground" onclick={() => {
        openResetPasswordPopup = true
        openLoginPopup = false
        }}>Reset Password</Button>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button type="submit">Login</Button>
    </Dialog.Footer>
  </form>
</Dialog.Content>
</Dialog.Root>

<!-- Register Popup -->
<Dialog.Root bind:open={openRegisterPopup}>
<Dialog.Content>
  <form onsubmit={register}>
    <Dialog.Header>

      <Label for="name">Name</Label>
      <Input id="name" name="name" type="text" bind:value={formData.name} />

      <Label for="email">E-Mail</Label>
      <Input id="email" name="email" type="email" bind:value={formData.email} />

      <Label for="password">Password</Label>
      <Input id="password" name="password" type="password" bind:value={formData.password} />

      <Label for="passwordConfirm">Confirm Password</Label>
      <Input id="passwordConfirm" name="passwordConfirm" type="password" bind:value={formData.passwordConfirm} />

      {#if formData.registerErrorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to Register</Alert.Title>
          <Alert.Description>
            {formData.registerErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button type="submit">Register</Button>
    </Dialog.Footer>
  </form>
</Dialog.Content>
</Dialog.Root>

<!-- Reset Password Popup -->
<Dialog.Root bind:open={openResetPasswordPopup}>
<Dialog.Content>
  <form onsubmit={resetPassword}>
    <Dialog.Header>

      <Label for="email">E-Mail</Label>
      <Input id="email" name="email" type="email" bind:value={formData.email} />

      {#if resetPasswordRequested}
        <Alert.Root>
          <HugeiconsIcon icon={InformationCircleIcon} />
          <Alert.Title>Password reset requested</Alert.Title>
          <Alert.Description>
            Please chek your E-Mail Inbox.
          </Alert.Description>
        </Alert.Root>
      {/if}
      
      {#if formData.resetPasswordErrorMessage}
        <Alert.Root variant="destructive">
          <HugeiconsIcon icon={Alert01Icon} />
          <Alert.Title>Unable to reset password</Alert.Title>
          <Alert.Description>
            {formData.resetPasswordErrorMessage}
          </Alert.Description>
        </Alert.Root>
      {/if}

    </Dialog.Header>
    <Dialog.Footer>
      <Dialog.Close>Cancel</Dialog.Close>
      <Button type="submit">Reset password</Button>
    </Dialog.Footer>
  </form>
</Dialog.Content>
</Dialog.Root>