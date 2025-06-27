import { render, fireEvent } from '@testing-library/react'
import { screen } from '@testing-library/dom'
import { vi } from 'vitest'

import Select from '@/components/Select'

test('renders dropdown and selects value', () => {
  const onChangeMock = vi.fn()

  render(
    <Select
      id="test-select"
      label="Test Select"
      options={[{ label: '456', value: '456' }]}
      value="123"
      onChange={onChangeMock}
    />
  )

  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '456' } })

  expect(onChangeMock).toHaveBeenCalledWith('456')
})
