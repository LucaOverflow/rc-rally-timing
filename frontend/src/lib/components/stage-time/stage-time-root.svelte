<script lang="ts">
  import * as Table from "$lib/components/ui/table"

  let { children } = $props();
</script>

<!-- Practice Table -->
<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head>Position</Table.Head>
      <Table.Head>Driver</Table.Head>
      <Table.Head>Best Time</Table.Head>
      <Table.Head>Last Time</Table.Head>
      <Table.Head>Diff to Next</Table.Head>
      <Table.Head>Diff to 1st</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {@render children?.()}
  </Table.Body>
</Table.Root>