<script setup lang="ts">
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import FloatLabel from 'primevue/floatlabel'
import type { Item, ItemForm } from '../types/item'

const props = defineProps<{
  visible: boolean
  item: Item | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  save: [item: ItemForm]
}>()

const form = ref<ItemForm>({
  code: '',
  name: '',
  stock: 0,
  location: '',
})

const errors = ref<Record<string, string>>({})

watch(
  () => props.visible,
  (val) => {
    if (val && props.item) {
      form.value = {
        code: props.item.code,
        name: props.item.name,
        stock: props.item.stock,
        location: props.item.location,
      }
    } else if (val) {
      form.value = { code: '', name: '', stock: 0, location: '' }
    }
    errors.value = {}
  }
)

const validate = (): boolean => {
  errors.value = {}

  if (!form.value.code.trim()) {
    errors.value.code = 'Code is required'
  }
  if (!form.value.name.trim()) {
    errors.value.name = 'Name is required'
  }
  if (form.value.stock < 0) {
    errors.value.stock = 'Stock cannot be negative'
  }

  return Object.keys(errors.value).length === 0
}

const handleSave = () => {
  if (validate()) {
    emit('save', { ...form.value })
  }
}

const handleClose = () => {
  emit('update:visible', false)
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="item ? 'Edit Item' : 'Create Item'"
    modal
    :style="{ width: '28rem' }"
    @update:visible="handleClose"
  >
    <div class="flex flex-col gap-4">
      <div>
        <FloatLabel variant="on">
          <InputText
            id="code"
            v-model="form.code"
            :invalid="!!errors.code"
            class="w-full"
          />
          <label for="code">Code *</label>
        </FloatLabel>
        <small v-if="errors.code" class="text-red-500">{{ errors.code }}</small>
      </div>

      <div>
        <FloatLabel variant="on">
          <InputText
            id="name"
            v-model="form.name"
            :invalid="!!errors.name"
            class="w-full"
          />
          <label for="name">Name *</label>
        </FloatLabel>
        <small v-if="errors.name" class="text-red-500">{{ errors.name }}</small>
      </div>

      <div>
        <FloatLabel variant="on">
          <InputNumber
            id="stock"
            v-model="form.stock"
            :invalid="!!errors.stock"
            class="w-full"
          />
          <label for="stock">Stock</label>
        </FloatLabel>
        <small v-if="errors.stock" class="text-red-500">{{ errors.stock }}</small>
      </div>

      <div>
        <FloatLabel variant="on">
          <InputText
            id="location"
            v-model="form.location"
            class="w-full"
          />
          <label for="location">Location</label>
        </FloatLabel>
      </div>
    </div>

    <template #footer>
      <Button label="Cancel" severity="secondary" text @click="handleClose" />
      <Button label="Save" icon="pi pi-check" @click="handleSave" />
    </template>
  </Dialog>
</template>
