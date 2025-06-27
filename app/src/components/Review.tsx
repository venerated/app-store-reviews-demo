import { format } from 'date-fns'

import type { IReview } from '@/types/index'

import styles from './review.module.scss'

export default function Review({ data }: { data: IReview }) {
  const formattedDate = format(data?.updated ?? '', 'PPP p')
  return (
    <div className={styles.wrap}>
      <div className={styles.header}>
        <div>{data.author}</div>
        <div className={styles.date}>{formattedDate}</div>
      </div>
      <div className={styles.meta}>
        <div>{data.rating}</div>
      </div>
      <div className={styles.content}>{data.content}</div>
    </div>
  )
}
