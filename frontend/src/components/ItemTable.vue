<script setup lang="ts">
import { ref, watch } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import type { Item } from '../types/item'

const props = defineProps<{
  items: Item[]
  loading: boolean
  totalRecords: number
  currentPage: number
  currentLimit: number
  search: string
}>()

const emit = defineEmits<{
  edit: [item: Item]
  delete: [item: Item]
  'page-change': [page: number, limit: number]
  search: [value: string]
}>()

const searchValue = ref(props.search)

let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchValue, (val) => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    emit('search', val)
  }, 300)
})

watch(() => props.search, (val) => {
  searchValue.value = val
})

const onPage = (event: any) => {
  const page = Math.floor(event.first / event.rows) + 1
  const limit = event.rows
  emit('page-change', page, limit)
}
</script>

<template>
  <div>
    <div class="mb-4">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="searchValue"
          placeholder="Search by name or code..."
          class="w-full md:w-80"
        />
      </IconField>
    </div>

    <DataTable
      :value="items"
      :loading="loading"
      lazy
      paginator
      :first="(currentPage - 1) * currentLimit"
      :rows="currentLimit"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      :totalRecords="totalRecords"
      stripedRows
      @page="onPage"
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
  </div>
</template>
