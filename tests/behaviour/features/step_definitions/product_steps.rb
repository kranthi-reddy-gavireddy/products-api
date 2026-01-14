require 'httparty'

BASE_URL = ENV['BASE_URL'] || 'http://localhost:8080'
Prouct_ID=''

When ('I request for healthcheck') do
  @response = HTTParty.get("#{BASE_URL}/health")
end

Then('I should receive {int} status with json payload')  do |status_code|
  expect(@response.code).to eq(status_code)
  expect(@response.parsed_response['status']).to eq('up')
end

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
  expect(@response.parsed_response).to have_key('id')
  Prouct_ID = @response.parsed_response['id']
  puts "Created Product ID: #{Prouct_ID}"
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

When('I publish Message to SQS for OrderCreatedTopic') do
    # curl -X POST "http://localhost:4566/000000000000/OrderCreatedTopic" \
    # -H "Content-Type: application/x-www-form-urlencoded" \
    # --data-urlencode "Action=SendMessage" \
    # --data-urlencode "MessageBody={\"product_id\":\"UN-20260114023132\",\"quantity\":22}" \
    # --data-urlencode "Version=2012-11-05"
    message_body = {
      product_id: Prouct_ID,
      quantity: 3
    }.to_json
    @response = HTTParty.post("http://localstack:4566/000000000000/OrderCreatedTopic",
                              body: {
                                "Action" => "SendMessage",
                                "MessageBody" => message_body,
                                "Version" => "2012-11-05"
                              },
                              headers: { 'Content-Type' => 'application/x-www-form-urlencoded' })
end

Then('the message should be published successfully') do
    expect(@response.code).to eq(200)
end

When('I request to the Product that is created') do
  @response = HTTParty.get("#{BASE_URL}/products/#{Prouct_ID}")
end

Then('I should receive the details of the Product with status {int}') do |status_code|
  expect(@response.code).to eq(status_code)
  expect(@response.parsed_response).to have_key('id')
  expect(@response.parsed_response['id']).to eq(Prouct_ID)
end

When('I request to delete the Product that is created') do
  @response = HTTParty.delete("#{BASE_URL}/products/#{Prouct_ID}")
end

Then('the Product should be deleted with status {int}') do |status_code|
  expect(@response.code).to eq(status_code)
  #return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": `Product deleted successfully %s`, "id": id})
  expect(@response.parsed_response['message']).to include('Product deleted successfully')
  expect(@response.parsed_response['message']).to include(Prouct_ID)
end