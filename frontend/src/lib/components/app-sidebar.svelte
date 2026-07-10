<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index"
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index"
  import {
    Calendar03Icon,
    CalendarClockIcon,
    ChevronUpIcon,
    InternetAntenna02Icon,
    RankingIcon,
    StopWatchIcon,
    Timer02Icon } from '@hugeicons/core-free-icons';
  import { HugeiconsIcon, type IconSvgElement } from '@hugeicons/svelte';

  interface MenuItem {
    title: string,
    url: string,
    icon: IconSvgElement
  }

  const driverItems: MenuItem[] = [
    {
      title: "Events",
      url: "#",
      icon: Calendar03Icon
    },
    {
      title: "Transponder",
      url: "#",
      icon: InternetAntenna02Icon
    },
    {
      title: "Decoder",
      url: "#",
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
  const isLoggedIn = false // TODO Make dynamic
</script>

<Sidebar.Root>

  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupLabel>Personal</Sidebar.GroupLabel>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          {#each driverItems as item (item.title)}
            <Sidebar.MenuItem>
              <Sidebar.MenuButton>
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

    {#if isEventSelected}
      <Sidebar.Group>
        <Sidebar.GroupLabel>Event</Sidebar.GroupLabel> <!-- TODO Add Event Name. Always show currently running Event as default. -->
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each eventItems as item (item.title)}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton>
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
                    Username <!-- TODO Add real username -->
                    <HugeiconsIcon icon={ChevronUpIcon} class="ms-auto" />
                  </Sidebar.MenuButton>
                {/snippet}
              </DropdownMenu.Trigger>
              <DropdownMenu.Content
                side="top"
                class="w-(--bits-dropdown-menu-anchor-width)"
              >
                <DropdownMenu.Item>
                  <span>Account</span>
                </DropdownMenu.Item>
                <DropdownMenu.Item>
                  <span>Sign out</span>
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          {:else}
            <Sidebar.Menu>
              <Sidebar.MenuItem>
                <Sidebar.MenuButton>
                  {#snippet child({ props })}
                    <a href="#" {...props}> <!-- TODO Open Log in Popup -->
                      <span>Log in</span>
                    </a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            </Sidebar.Menu>
          {/if}
        </Sidebar.MenuItem>
      </Sidebar.Menu>
  </Sidebar.Footer>

</Sidebar.Root>