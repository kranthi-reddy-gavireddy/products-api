require 'httparty'

BASE_URL = ENV['BASE_URL'] || 'http://localhost:8080'

When('I request all products') do
  @response = HTTParty.get("#{BASE_URL}/products")
end

Then('I should receive a list of products') do
  expect(@response.code).to eq(200)
  expect(@response.parsed_response).to be_an(Array)
end

When('I create a product with name {string}, price {float}, seller_id {string}, quantity {int}') do |name, price, seller_id, quantity|
  @product_data = {
    name: name,
    price: price,
    seller_id: seller_id,
    quantity: quantity
  }
  @response = HTTParty.post("#{BASE_URL}/products",
                            body: @product_data.to_json,
                            headers: { 'Content-Type' => 'application/json' })
end

Then('the product should be created successfully') do
  expect(@response.code).to eq(201)
end

Then('I should receive the product details') do
  expect(@response.parsed_response).to have_key('id')
  expect(@response.parsed_response['name']).to eq(@product_data[:name])
  expect(@response.parsed_response['price']).to eq(@product_data[:price])
end

When('I create a product with invalid data') do
  invalid_data = {
    name: '',
    price: -1
  }
  @response = HTTParty.post("#{BASE_URL}/products",
                            body: invalid_data.to_json,
                            headers: { 'Content-Type' => 'application/json' })
end

Then('I should receive an error') do
  expect(@response.code).to eq(400)
end