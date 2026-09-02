import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { TraceExplorer } from './TraceExplorer'
import { trace } from '../test/fixtures'

describe('TraceExplorer',()=>{it('filters trace outcomes and expands waterfalls',async()=>{const user=userEvent.setup();render(<TraceExplorer traces={[trace]}/>);expect(screen.getByText('Authorize')).toBeInTheDocument();await user.click(screen.getByRole('button',{name:/gateway/}));expect(screen.queryByText('Authorize')).not.toBeInTheDocument();await user.selectOptions(screen.getByLabelText('Trace status'),'ok');expect(screen.getByText('No traces match the active investigation filters.')).toBeInTheDocument()})})
