import axios from 'axios'
import type { Item, ItemForm, PaginatedResponse, ItemFilter } from '../types/item'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

export const itemApi = {
  async getAll(filter: ItemFilter = {}): Promise<PaginatedResponse<Item>> {
    const params = new URLSearchParams()
    if (filter.search) params.set('search', filter.search)
    if (filter.location) params.set('location', filter.location)
    if (filter.min_stock && filter.min_stock > 0) params.set('min_stock', String(filter.min_stock))
    if (filter.max_stock && filter.max_stock > 0) params.set('max_stock', String(filter.max_stock))
    if (filter.page && filter.page > 0) params.set('page', String(filter.page))
    if (filter.limit && filter.limit > 0) params.set('limit', String(filter.limit))

    const queryString = params.toString()
    const { data } = await api.get<PaginatedResponse<Item>>(`/items${queryString ? '?' + queryString : ''}`)
    return data
  },

  async getById(id: number): Promise<Item> {
    const { data } = await api.get<Item>(`/items/${id}`)
    return data
  },

  async create(item: ItemForm): Promise<Item> {
    const { data } = await api.post<Item>('/items', item)
    return data
  },

  async update(id: number, item: ItemForm): Promise<Item> {
    const { data } = await api.put<Item>(`/items/${id}`, item)
    return data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/items/${id}`)
  },

  async exportCsv(filter: ItemFilter = {}): Promise<Blob> {
    const params = new URLSearchParams()
    if (filter.search) params.set('search', filter.search)
    if (filter.location) params.set('location', filter.location)
    if (filter.min_stock && filter.min_stock > 0) params.set('min_stock', String(filter.min_stock))
    if (filter.max_stock && filter.max_stock > 0) params.set('max_stock', String(filter.max_stock))

    const queryString = params.toString()
    const { data } = await api.get(`/items/export${queryString ? '?' + queryString : ''}`, {
      responseType: 'blob',
    })
    return data
  },
}
