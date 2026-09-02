import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ErrorBudgets } from './ErrorBudgets'
import { slo } from '../test/fixtures'

describe('ErrorBudgets', () => {
  it('shows target attainment and remaining budget', () => {
    render(<ErrorBudgets slos={[slo]} />)

    expect(screen.getByText('Payment success')).toBeInTheDocument()
    expect(screen.getByText(/99\.00% target/)).toBeInTheDocument()
    expect(screen.getByText(/50% remaining/)).toBeInTheDocument()
  })
})
