export interface Item {
  id: number
  code: string
  name: string
  stock: number
  location: string
}

export interface ItemForm {
  code: string
  name: string
  stock: number
  location: string
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface ItemFilter {
  search?: string
  location?: string
  min_stock?: number
  max_stock?: number
  page?: number
  limit?: number
}

export interface ValidationError {
  field: string
  message: string
}

export interface ValidationErrors {
  errors: ValidationError[]
}

export interface ApiError {
  error: {
    code: string
    message: string
  }
}
