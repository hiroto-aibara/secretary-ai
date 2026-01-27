export function generateListId(name: string): string {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

export function generateUniqueListId(
  name: string,
  existingIds: Set<string>,
): string {
  const baseId = generateListId(name)
  if (!existingIds.has(baseId)) {
    return baseId
  }

  let counter = 1
  let newId = `${baseId}-${counter}`
  while (existingIds.has(newId)) {
    counter++
    newId = `${baseId}-${counter}`
  }
  return newId
}
