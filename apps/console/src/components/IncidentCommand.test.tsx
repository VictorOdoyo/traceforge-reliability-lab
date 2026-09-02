import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { IncidentCommand } from './IncidentCommand'
import { incident } from '../test/fixtures'

describe('IncidentCommand',()=>{it('advances an incident along the response lifecycle',async()=>{const onTransition=vi.fn();render(<IncidentCommand incidents={[incident]} onTransition={onTransition}/>);expect(screen.getByText('Fast burn detected')).toBeInTheDocument();await userEvent.click(screen.getByRole('button',{name:/Move to identified/}));expect(onTransition).toHaveBeenCalledWith('inc-1','identified')})})
