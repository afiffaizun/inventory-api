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

export interface ApiError {
  error: {
    code: string
    message: string
  }
}
