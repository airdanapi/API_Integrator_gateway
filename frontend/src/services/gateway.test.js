import { describe, expect, it, vi } from 'vitest'
import {
  sendGatewayLogistics,
  sendGatewayMarketplace,
  sendGatewayPayment,
  sendGatewaySmartBank,
  sendGatewaySupplier,
} from './gateway'

describe('gateway service', () => {
  it('sends payment request', async () => {
    const payload = {
      from_app: 'Marketplace',
      from_user: 'user1',
      to_user: 'user2',
      amount: 10000,
      service_type: 'payment',
    }
    const mockClient = {
      post: vi.fn().mockResolvedValue({
        data: {
          status: 'success',
          data: { transaction_id: 'gw-payment-1', forwarded: true, upstream: 'smartbank' },
        },
      }),
    }

    const result = await sendGatewayPayment(mockClient, payload)

    expect(mockClient.post).toHaveBeenCalledWith('/gateway/payment', payload)
    expect(result.transaction_id).toBe('gw-payment-1')
    expect(result.forwarded).toBe(true)
  })

  it('sends smartbank request', async () => {
    const payload = { action: 'check_balance', payload: { account: '123' } }
    const mockClient = {
      post: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { forwarded: true, upstream: 'smartbank' } },
      }),
    }

    const result = await sendGatewaySmartBank(mockClient, payload)

    expect(mockClient.post).toHaveBeenCalledWith('/gateway/smartbank', payload)
    expect(result.forwarded).toBe(true)
  })

  it('sends marketplace request', async () => {
    const payload = { action: 'get_order', payload: { order_id: 'ORD-1' } }
    const mockClient = {
      post: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { forwarded: true, upstream: 'marketplace' } },
      }),
    }

    const result = await sendGatewayMarketplace(mockClient, payload)

    expect(mockClient.post).toHaveBeenCalledWith('/gateway/marketplace', payload)
    expect(result.forwarded).toBe(true)
  })

  it('sends logistics request', async () => {
    const payload = {
      order_id: 'ORD-1',
      address: 'Jl. Merdeka No. 1',
      distance: 10,
      shipping_type: 'express',
    }
    const mockClient = {
      post: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { delivery_id: 'DEL-1', forwarded: true, upstream: 'logistikit' } },
      }),
    }

    const result = await sendGatewayLogistics(mockClient, payload)

    expect(mockClient.post).toHaveBeenCalledWith('/gateway/logistics', payload)
    expect(result.delivery_id).toBe('DEL-1')
  })

  it('sends supplier request', async () => {
    const payload = { supplier_id: 'SUP-1', material: 'Beras', qty: 50, total_cost: 250000 }
    const mockClient = {
      post: vi.fn().mockResolvedValue({
        data: { status: 'success', data: { order_id: 'SUP-001', forwarded: true, upstream: 'supplierhub' } },
      }),
    }

    const result = await sendGatewaySupplier(mockClient, payload)

    expect(mockClient.post).toHaveBeenCalledWith('/gateway/supplier', payload)
    expect(result.order_id).toBe('SUP-001')
  })
})
