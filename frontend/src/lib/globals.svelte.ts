import type { RecordModel } from 'pocketbase'

export const globals: {
  activeEvent: RecordModel | undefined
} = $state({
  activeEvent: undefined
})