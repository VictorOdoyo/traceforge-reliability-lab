import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { AlertQueue } from './AlertQueue'
import { alert } from '../test/fixtures'

describe('AlertQueue',()=>{it('allows an operator to claim a firing alert',async()=>{const onAcknowledge=vi.fn();render(<AlertQueue alerts={[alert]} onAcknowledge={onAcknowledge}/>);await userEvent.click(screen.getByRole('button',{name:'Acknowledge'}));expect(onAcknowledge).toHaveBeenCalledWith('alert-1')})})
