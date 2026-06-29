<script setup lang="ts">
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import type { Item } from '../types/item'

defineProps<{
  items: Item[]
  loading: boolean
}>()

const emit = defineEmits<{
  edit: [item: Item]
  delete: [item: Item]
}>()
</script>

<template>
  <DataTable
    :value="items"
    :loading="loading"
    paginator
    :rows="10"
    :rowsPerPageOptions="[5, 10, 20, 50]"
    stripedRows
    class="rounded-lg shadow"
  >
    <Column field="code" header="Code" sortable style="min-width: 120px" />
    <Column field="name" header="Name" sortable style="min-width: 200px" />
    <Column field="stock" header="Stock" sortable style="min-width: 100px" />
    <Column field="location" header="Location" sortable style="min-width: 150px" />
    <Column header="Actions" style="min-width: 150px">
      <template #body="{ data }">
        <div class="flex gap-2">
          <Button
            icon="pi pi-pencil"
            severity="info"
            text
            rounded
            @click="emit('edit', data)"
            v-tooltip.top="'Edit'"
          />
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            @click="emit('delete', data)"
            v-tooltip.top="'Delete'"
          />
        </div>
      </template>
    </Column>
    <template #empty>
      <div class="text-center py-8 text-gray-500">
        <i class="pi pi-inbox text-4xl mb-2"></i>
        <p>No items found</p>
      </div>
    </template>
  </DataTable>
</template>
