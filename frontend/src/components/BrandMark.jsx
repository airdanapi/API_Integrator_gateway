function BrandMark({ compact = false }) {
  return (
    <span className="inline-flex items-center gap-3">
      <span
        className={`${compact ? 'h-9 w-9' : 'h-11 w-11'} grid shrink-0 place-items-center rounded-xl bg-blue-700 text-white shadow-lg shadow-blue-700/20`}
        aria-hidden="true"
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          className={compact ? 'h-5 w-5' : 'h-6 w-6'}
        >
          <path
            d="M7 7.5h4.25M7 16.5h4.25M12.75 7.5H17M12.75 16.5H17M9.2 7.5a2.2 2.2 0 1 1-4.4 0 2.2 2.2 0 0 1 4.4 0ZM19.2 7.5a2.2 2.2 0 1 1-4.4 0 2.2 2.2 0 0 1 4.4 0ZM9.2 16.5a2.2 2.2 0 1 1-4.4 0 2.2 2.2 0 0 1 4.4 0ZM19.2 16.5a2.2 2.2 0 1 1-4.4 0 2.2 2.2 0 0 1 4.4 0Z"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
          />
        </svg>
      </span>
      <span className="leading-tight">
        <span className="block text-sm font-extrabold tracking-tight text-slate-950 sm:text-base">
          API Integrator
        </span>
        <span className="block text-[0.65rem] font-semibold uppercase tracking-[0.18em] text-blue-700">
          Gateway UMKM
        </span>
      </span>
    </span>
  )
}

export default BrandMark
