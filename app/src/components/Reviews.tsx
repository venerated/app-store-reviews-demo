import { getUnixTime, parseISO } from 'date-fns'
import { useEffect, useMemo, useState } from 'react'

import Review from '@/components/Review'
import Select from '@/components/Select'

import type { IReview } from '@/types/index'

import styles from './reviews.module.scss'

const TIMEFRAME = 48

export default function Reviews() {
  // App
  const [appId, setAppId] = useState<string | null>(null)

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

  // Reviews
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [reviews, setReviews] = useState<IReview[] | null>(null)

  const handleError = (msg: string) => {
    console.error(msg)
    setError(msg)
  }

  useEffect(() => {
    if (!appId) return

    // Fetch reviews from backend
    const fetchReviews = async () => {
      setLoading(true)
      try {
        const response = await fetch(
          `${import.meta.env.VITE_API_URL}/reviews?id=${appId}`
        )
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        const result = await response.json()
        setReviews(result)
      } catch (err: unknown) {
        if (err instanceof Error) {
          handleError(`An error occured: ${err.message}`)
        } else {
          handleError(`An unexpected error occured: ${err}`)
        }
      } finally {
        setLoading(false)
      }
    }

    fetchReviews()
  }, [appId])

  const filteredReviews = useMemo(() => {
    const currentUnixTime = Math.floor(Date.now() / 1000)
    // Unix timestamp for set timeframe
    const earliestAllowedTimestamp = currentUnixTime - 60 * 60 * TIMEFRAME
    return reviews
      ?.filter((review) => {
        const reviewTimestamp = review.updated
        const unixTimestamp = getUnixTime(parseISO(reviewTimestamp))
        return unixTimestamp > earliestAllowedTimestamp
      })
      .sort((a, b) => {
        const aTime = getUnixTime(parseISO(a.updated))
        const bTime = getUnixTime(parseISO(b.updated))
        // Sort descending
        return bTime - aTime
      })
  }, [reviews])

  return (
    <div className={styles.wrap}>
      <h2>Reviews</h2>

      <div className={styles.filters}>
        <Select
          id="app-options"
          label="App to Display Reviews For"
          options={appOptions}
          placeholder="Choose an App"
          selected={appId}
          onChange={setAppId}
        />
      </div>

      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div>{error}</div>
      ) : (
        <div className={styles.reviews}>
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
