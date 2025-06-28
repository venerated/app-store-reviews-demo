import { format } from 'date-fns'

import type { IReview } from '@/types/index'

import styles from './review.module.scss'

export default function Review({ data }: { data: IReview }) {
  const formattedDate = format(data.updated ?? '', 'PPP p')
  const stars = [...Array(Number(data.rating ?? 0)).keys()]

  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <div className={styles.author}>{data.author}</div>
        <div className={styles.date}>{formattedDate}</div>
      </div>
      <div className={styles.meta}>
        <div className={styles.stars}>
          {stars.length
            ? stars.map((_, index) => <span key={index}>⭐️</span>)
            : null}
        </div>
      </div>
      <div className={styles.content}>{data.content}</div>
    </div>
  )
}
