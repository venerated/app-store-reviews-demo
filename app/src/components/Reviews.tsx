import { getUnixTime, parseISO } from 'date-fns'
import { useEffect, useMemo, useState } from 'react'

import Review from '@/components/Review'
import Select from '@/components/Select'

import type { IReview } from '@/types/index'

import styles from './reviews.module.scss'

const TIMEFRAME = 48

// This component allows the user to select an app from a dropdown menu,
// fetches reviews for that app from the backend, filters them by a recent
// timeframe, and displays them sorted by date.
export default function Reviews() {
  const [appId, setAppId] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [reviews, setReviews] = useState<IReview[] | null>(null)

  // Hardcoded list of popular apps to choose from, each with a label and its corresponding App Store ID
  const appOptions = [
    {
      label: 'ChatGPT (6448311069)',
      value: '6448311069',
    },
    {
      label: 'Threads (6446901002)',
      value: '6446901002',
    },
    {
      label: 'Peacock TV: Stream TV & Movies (1508186374)',
      value: '1508186374',
    },
    {
      label: 'Fortnite (1181774280)',
      value: '1181774280',
    },
    {
      label: 'Google (284815942)',
      value: '284815942',
    },
    {
      label: 'DAZN: Stream Live Sports (1129523589)',
      value: '1129523589',
    },
    {
      label: 'Netflix (363590051)',
      value: '363590051',
    },
    {
      label: 'Astra - Life Advice (6448947021)',
      value: '6448947021',
    },
    {
      label: 'Starla - Call the Universe (6448913133)',
      value: '6448913133',
    },
    {
      label: 'Max: Stream HBO, TV, & Movies (1660221872)',
      value: '1660221872',
    },
  ]

  const handleError = (msg: string) => {
    console.error(msg)
    setError(msg)
  }

  // Fetch reviews from backend whenever the selected app ID changes
  useEffect(() => {
    if (!appId) return

    const fetchReviews = async () => {
      setLoading(true)
      try {
        const response = await fetch(
          `${import.meta.env.VITE_API_URL}/reviews?id=${appId}`
        )
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${String(response.status)}`)
        }
        const result = (await response.json()) as IReview[]
        setReviews(result)
      } catch (err: unknown) {
        if (err instanceof Error) {
          handleError(`An error occured: ${err.message}`)
        } else {
          handleError(
            `An unexpected error occured: ${
              typeof err === 'string' ? err : JSON.stringify(err)
            }`
          )
        }
      } finally {
        setLoading(false)
      }
    }

    void fetchReviews()
  }, [appId])

  // Memoized computation of filtered and sorted reviews, only includes reviews
  // newer than the configured timeframe (e.g. 48 hours).
  // Sorted from most recent to oldest
  const filteredReviews = useMemo(() => {
    const currentUnixTime = Math.floor(Date.now() / 1000)
    const earliestAllowedTimestamp = currentUnixTime - 60 * 60 * TIMEFRAME
    return reviews
      ?.filter((review) => {
        const reviewTimestamp = review.updated
        // Fallback to current timestamp if field is blank from RSS
        const unixTimestamp = reviewTimestamp
          ? getUnixTime(parseISO(reviewTimestamp))
          : currentUnixTime
        return unixTimestamp > earliestAllowedTimestamp
      })
      .sort((a, b) => {
        // Fallback to current timestamp if field is blank from RSS
        const aTime = a.updated
          ? getUnixTime(parseISO(a.updated))
          : currentUnixTime
        const bTime = b.updated
          ? getUnixTime(parseISO(b.updated))
          : currentUnixTime
        return bTime - aTime
      })
  }, [reviews])

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <h2>App Store Reviews</h2>
      </div>
      <div className={styles.filters}>
        <Select
          id="app-options"
          label="App to Display Reviews For"
          options={appOptions}
          placeholder="Choose an App"
          value={appId}
          onChange={setAppId}
        />
      </div>

      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div>{error}</div>
      ) : (
        <div className={styles.reviews} role="group" aria-label="Reviews">
          {filteredReviews?.length ? (
            filteredReviews.map((review) => (
              <Review key={review.id} data={review} />
            ))
          ) : !filteredReviews?.length && appId ? (
            <div>No Reviews in the last {TIMEFRAME} hours</div>
          ) : null}
        </div>
      )}
    </div>
  )
}
