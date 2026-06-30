<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Button from 'primevue/button'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import ItemTable from '../components/ItemTable.vue'
import ItemForm from '../components/ItemForm.vue'
import { itemApi } from '../api/item'
import type { Item, ItemForm as ItemFormType, ItemFilter } from '../types/item'

const toast = useToast()
const confirm = useConfirm()

const items = ref<Item[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingItem = ref<Item | null>(null)
const totalRecords = ref(0)
const currentPage = ref(1)
const currentLimit = ref(10)
const currentSearch = ref('')

const filter = ref<ItemFilter>({
  page: 1,
  limit: 10,
})

const fetchItems = async () => {
  loading.value = true
  try {
    const response = await itemApi.getAll(filter.value)
    items.value = response.data
    totalRecords.value = response.total
    currentPage.value = response.page
    currentLimit.value = response.limit
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to fetch items',
      life: 3000,
    })
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  editingItem.value = null
  dialogVisible.value = true
}

const handleEdit = (item: Item) => {
  editingItem.value = item
  dialogVisible.value = true
}

const handleSave = async (form: ItemFormType) => {
  try {
    if (editingItem.value) {
      await itemApi.update(editingItem.value.id, form)
      toast.add({
        severity: 'success',
        summary: 'Updated',
        detail: 'Item updated successfully',
        life: 3000,
      })
    } else {
      await itemApi.create(form)
      toast.add({
        severity: 'success',
        summary: 'Created',
        detail: 'Item created successfully',
        life: 3000,
      })
    }
    dialogVisible.value = false
    await fetchItems()
  } catch (err: any) {
    const message = err?.response?.data?.errors
      ? err.response.data.errors.map((e: any) => e.message).join(', ')
      : 'Failed to save item'
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: message,
      life: 4000,
    })
  }
}

const handleDelete = (item: Item) => {
  confirm.require({
    message: `Are you sure you want to delete "${item.name}"?`,
    header: 'Delete Confirmation',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    rejectClass: 'p-button-secondary',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await itemApi.delete(item.id)
        toast.add({
          severity: 'success',
          summary: 'Deleted',
          detail: 'Item deleted successfully',
          life: 3000,
        })
        await fetchItems()
      } catch {
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to delete item',
          life: 3000,
        })
      }
    },
  })
}

const handlePageChange = (page: number, limit: number) => {
  filter.value.page = page
  filter.value.limit = limit
  fetchItems()
}

const handleSearch = (search: string) => {
  currentSearch.value = search
  filter.value.search = search
  filter.value.page = 1
  fetchItems()
}

const handleExport = async () => {
  try {
    const blob = await itemApi.exportCsv(filter.value)
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'items.csv'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    toast.add({
      severity: 'success',
      summary: 'Exported',
      detail: 'CSV downloaded successfully',
      life: 3000,
    })
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to export items',
      life: 3000,
    })
  }
}

onMounted(() => {
  fetchItems()
})
</script>

<template>
  <div>
    <div class="max-w-7xl mx-auto px-4 py-8">
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-800">Inventory Items</h1>
          <p class="text-gray-500 mt-1">Manage your inventory items</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Export CSV"
            icon="pi pi-download"
            severity="secondary"
            outlined
            @click="handleExport"
          />
          <Button
            label="Add Item"
            icon="pi pi-plus"
            @click="handleCreate"
          />
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-4">
        <ItemTable
          :items="items"
          :loading="loading"
          :total-records="totalRecords"
          :current-page="currentPage"
          :current-limit="currentLimit"
          :search="currentSearch"
          @edit="handleEdit"
          @delete="handleDelete"
          @page-change="handlePageChange"
          @search="handleSearch"
        />
      </div>
    </div>

    <ItemForm
      v-model:visible="dialogVisible"
      :item="editingItem"
      @save="handleSave"
    />
  </div>
</template>
