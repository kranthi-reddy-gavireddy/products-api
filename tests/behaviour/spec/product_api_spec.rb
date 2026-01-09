require 'rspec'
require 'httparty'

RSpec.describe 'Products API' do
  let(:base_url) { 'http://localhost:8080' }

  describe 'GET /products' do
    it 'returns a list of products' do
      response = HTTParty.get("#{base_url}/products")
      expect(response.code).to eq(200)
      expect(response.parsed_response).to be_an(Array)
    end
  end

  describe 'POST /products' do
    it 'creates a new product' do
      product_data = {
        name: 'Test Product',
        price: 10.99,
        seller_id: 'seller123',
        quantity: 5
      }

      response = HTTParty.post("#{base_url}/products",
                               body: product_data.to_json,
                               headers: { 'Content-Type' => 'application/json' })

      expect(response.code).to eq(201)
      expect(response.parsed_response).to have_key('id')
      expect(response.parsed_response['name']).to eq('Test Product')
    end

    it 'returns error for invalid data' do
      invalid_data = {
        name: '',
        price: -1
      }

      response = HTTParty.post("#{base_url}/products",
                               body: invalid_data.to_json,
                               headers: { 'Content-Type' => 'application/json' })

      expect(response.code).to eq(400)
    end
  end
end