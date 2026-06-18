function SectionHeading({ eyebrow, title, description, align = 'left' }) {
  const alignment =
    align === 'center'
      ? 'mx-auto max-w-3xl text-center'
      : 'max-w-2xl text-left'

  return (
    <div className={alignment}>
      <p className="mb-3 text-xs font-bold uppercase tracking-[0.22em] text-blue-700">
        {eyebrow}
      </p>
      <h2 className="text-balance text-3xl font-extrabold tracking-tight text-slate-950 sm:text-4xl">
        {title}
      </h2>
      {description ? (
        <p className="mt-4 text-base leading-7 text-slate-600 sm:text-lg">
          {description}
        </p>
      ) : null}
    </div>
  )
}

export default SectionHeading
