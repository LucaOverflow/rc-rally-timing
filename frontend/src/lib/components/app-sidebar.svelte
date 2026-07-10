<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index"
  import {
    Calendar03Icon,
    CalendarClockIcon,
    InternetAntenna02Icon,
    RankingIcon,
    StopWatchIcon,
    Timer02Icon,
    UserAccountIcon } from '@hugeicons/core-free-icons';
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
    },
    {
      title: "Personal Data",
      url: "#",
      icon: UserAccountIcon
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
</script>

<Sidebar.Root>

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

</Sidebar.Root>