import styles from './select.module.scss'

const Select = ({
  id,
  label,
  options,
  placeholder,
  selected,
  onChange,
}: {
  id: string
  label: string
  options: { label: string; value: string }[]
  placeholder: string
  selected: string | null
  onChange: (val: string) => void
}) => {
  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    onChange(event.target.value)
  }

  return (
    <div className={styles.wrap}>
      <label htmlFor={id}>{label}</label>
      <select id={id} value={selected ?? ''} onChange={handleChange}>
        <option value="" disabled>
          {placeholder ?? 'Make a Selection'}
        </option>
        {options?.length
          ? options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))
          : null}
      </select>
    </div>
  )
}

export default Select
