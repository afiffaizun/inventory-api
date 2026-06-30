<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()
const sidebarCollapsed = ref(false)

const menuItems = [
  { label: 'Items', icon: 'pi pi-box', route: '/' },
  { label: 'Categories', icon: 'pi pi-tags', route: '/categories' },
  { label: 'Reports', icon: 'pi pi-chart-bar', route: '/reports' },
]

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex">
    <Toast />
    <ConfirmDialog />

    <!-- Sidebar -->
    <aside
      :class="[
        'bg-gray-900 text-white flex flex-col transition-all duration-300',
        sidebarCollapsed ? 'w-16' : 'w-64'
      ]"
    >
      <!-- Logo / Brand -->
      <div class="flex items-center gap-3 px-4 py-5 border-b border-gray-700">
        <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center flex-shrink-0">
          <i class="pi pi-box text-white text-sm"></i>
        </div>
        <transition name="fade">
          <span v-if="!sidebarCollapsed" class="font-bold text-lg whitespace-nowrap">
            Inventory
          </span>
        </transition>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 py-4">
        <ul class="space-y-1 px-2">
          <li v-for="item in menuItems" :key="item.route">
            <button
              @click="navigateTo(item.route)"
              :class="[
                'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors',
                route.path === item.route
                  ? 'bg-gray-700 text-white'
                  : 'text-gray-400 hover:bg-gray-800 hover:text-white'
              ]"
            >
              <i :class="[item.icon, 'text-lg']"></i>
              <transition name="fade">
                <span v-if="!sidebarCollapsed" class="whitespace-nowrap">
                  {{ item.label }}
                </span>
              </transition>
            </button>
          </li>
        </ul>
      </nav>

      <!-- Collapse Toggle -->
      <div class="border-t border-gray-700 p-2">
        <Button
          :icon="sidebarCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"
          text
          rounded
          @click="toggleSidebar"
          class="w-full text-gray-400 hover:text-white"
        />
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 overflow-auto">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
