import {
  render,
  fireEvent,
  screen,
  waitFor,
  within,
} from '@testing-library/react'
import '@testing-library/jest-dom'
import Reviews from '@/components/Reviews'
import { vi } from 'vitest'

// Mock data representing a typical review returned from the API
const mockReviews = [
  {
    id: '1',
    updated: new Date().toISOString(),
    content: 'Great app!',
    user: 'Tester',
  },
]

// Before each test, mock the global fetch function to return a successful response
// with the mockReviews array. This ensures consistent test data and isolates tests
// from real network calls.
beforeEach(() => {
  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve(mockReviews),
    })
  ) as unknown as typeof fetch
})

/**
 * Test: Loads and displays reviews when an app is selected
 * - Renders the Reviews component
 * - Simulates selecting an app by changing the combobox value
 * - Waits for the review text to appear in the document, confirming that
 *   reviews were fetched and displayed correctly.
 */
it('loads and displays reviews when app is selected', async () => {
  render(<Reviews />)

  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  await waitFor(() => {
    expect(screen.getByText(/Great app!/i)).toBeInTheDocument()
  })
})

/**
 * Test: Shows no reviews message if API returns an empty array
 * - Overrides the fetch mock to return an empty array, simulating no reviews
 * - Renders the Reviews component and selects an app
 * - Asserts that the previously mocked review text does not appear
 */
it('shows no reviews message if API returns empty array', async () => {
  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve([]),
    })
  ) as unknown as typeof fetch

  render(<Reviews />)
  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  await waitFor(() => {
    expect(screen.queryByText(/Great app!/i)).not.toBeInTheDocument()
  })
})

/**
 * Test: Handles fetch errors gracefully
 * - Mocks fetch to reject with an error simulating a network failure
 * - Renders the component and selects an app
 * - Confirms that error message is displayed
 */
it('handles fetch errors gracefully', async () => {
  global.fetch = vi.fn(() =>
    Promise.reject(new Error('Fetch failed'))
  ) as unknown as typeof fetch

  render(<Reviews />)

  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  await waitFor(() => {
    expect(
      screen.getByText(/An error occured: Fetch failed/i)
    ).toBeInTheDocument()
  })
})

/**
 * Test: Does not show reviews older than 48 hours
 * - Creates a mock review with a timestamp older than 48 hours
 * - Mocks fetch to return this old review
 * - Renders the component and selects an app
 * - Asserts that the old review is filtered out and not displayed
 */
it('does not show reviews older than 48 hours', async () => {
  const oldDate = new Date(Date.now() - 1000 * 60 * 60 * 72).toISOString()

  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () =>
        Promise.resolve([
          { id: '2', updated: oldDate, content: 'Old review', user: 'OldUser' },
        ]),
    })
  ) as unknown as typeof fetch

  render(<Reviews />)
  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  await waitFor(() => {
    expect(screen.queryByText(/Old review/i)).not.toBeInTheDocument()
  })
})

/**
 * Test: Sorts reviews by updated date descending
 * - Creates two reviews with different timestamps (newer and older)
 * - Mocks fetch to return these reviews unordered
 * - Renders the component and selects an app
 * - Asserts that the reviews are displayed sorted with the newest first
 */
it('sorts reviews by updated date descending', async () => {
  const now = new Date()
  const older = new Date(now.getTime() - 1000 * 60 * 60).toISOString()
  const newer = new Date(now.getTime()).toISOString()

  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () =>
        Promise.resolve([
          { id: '1', updated: older, content: 'Older review', user: 'User1' },
          { id: '2', updated: newer, content: 'Newer review', user: 'User2' },
        ]),
    })
  ) as unknown as typeof fetch

  render(<Reviews />)
  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  await waitFor(() => {
    const reviewsGroup = screen.getByRole('group', { name: /reviews/i })
    const reviewElements = within(reviewsGroup).getAllByText(/review/i)
    expect(reviewElements[0]).toHaveTextContent('Newer review')
    expect(reviewElements[1]).toHaveTextContent('Older review')
  })
})

/**
 * Test: Shows loading indicator while fetching reviews
 * - Creates a fetch mock that resolves only when manually triggered
 * - Renders the component and selects an app
 * - Asserts that a loading indicator is shown before fetch resolves
 * - After resolving fetch, asserts that loading indicator disappears and reviews appear
 */
it('shows loading indicator while fetching reviews', async () => {
  // Placeholder function to resolve fetch later
  let resolveFetch: () => void = () => {}

  // Create a promise that will resolve when resolveFetch is called
  const fetchPromise = new Promise((resolve) => {
    resolveFetch = () =>
      resolve({
        ok: true,
        json: () => Promise.resolve(mockReviews),
      })
  })

  // Mock fetch to return the above promise, delaying resolution
  global.fetch = vi.fn(() => fetchPromise as any) as typeof fetch

  render(<Reviews />)
  const select = screen.getByRole('combobox')
  fireEvent.change(select, { target: { value: '6448311069' } })

  expect(screen.getByText(/loading/i)).toBeInTheDocument()

  resolveFetch()

  await waitFor(() => {
    expect(screen.queryByText(/loading/i)).not.toBeInTheDocument()
    expect(screen.getByText(/Great app!/i)).toBeInTheDocument()
  })
})

/**
 * Test: Does not fetch reviews until an app is selected
 * - Sets up a spy on fetch to track calls
 * - Renders the Reviews component without selecting an app
 * - Asserts that fetch was not called since no app was selected
 */
it('does not fetch reviews until an app is selected', async () => {
  const fetchSpy = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve(mockReviews),
    })
  )
  global.fetch = fetchSpy as unknown as typeof fetch

  render(<Reviews />)

  expect(fetchSpy).not.toHaveBeenCalled()
})

/**
 * Test: Handles invalid or missing app selection gracefully
 * - Mocks fetch to return mockReviews (though it should not be called)
 * - Renders the component without selecting any app
 * - Waits to confirm that no reviews are displayed since no app was selected
 */
it('handles invalid or missing app selection gracefully', async () => {
  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve(mockReviews),
    })
  ) as unknown as typeof fetch

  render(<Reviews />)

  await waitFor(() => {
    expect(screen.queryByText(/Great app!/i)).not.toBeInTheDocument()
  })
})
