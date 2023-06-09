export abstract class Entity<T> {
  protected readonly _id: string = null
  protected readonly _createdAt: Date = null

  protected readonly data: T

  constructor(data: T, id: string, createdAt?: Date) {
    this._id = id
    this.data = data
    this._createdAt = createdAt || new Date()
  }

  get id() {
    return this._id
  }

  get createdAt() {
    return this._createdAt
  }
}
