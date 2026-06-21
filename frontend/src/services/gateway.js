export async function sendGatewayPayment(apiClient, payload) {
  const response = await apiClient.post('/gateway/payment', payload)
  return response.data.data
}

export async function sendGatewaySmartBank(apiClient, payload) {
  const response = await apiClient.post('/gateway/smartbank', payload)
  return response.data.data
}

export async function sendGatewayMarketplace(apiClient, payload) {
  const response = await apiClient.post('/gateway/marketplace', payload)
  return response.data.data
}

export async function sendGatewayLogistics(apiClient, payload) {
  const response = await apiClient.post('/gateway/logistics', payload)
  return response.data.data
}

export async function sendGatewaySupplier(apiClient, payload) {
  const response = await apiClient.post('/gateway/supplier', payload)
  return response.data.data
}
