Feature: Product Management
  As a user of the products API
  I want to manage products
  So that I can create and retrieve product information

  Scenario: Trigger healthcheck point
            When I request for healthcheck
            Then I should receive 200 status with json payload

  Scenario: Create a new product
    When I create a product with name "Test Product", price 10.99, seller_id "seller123", quantity 5
    Then the product should be created successfully
    And I should receive the product details
  
  Scenario: Updating the product count by publishing OrderCreatedTopic 
    When I publish Message to SQS for OrderCreatedTopic
    Then the message should be published successfully

  Scenario: Retrieve all products
    When I request all products
    Then I should receive a list of products
  
  Scenario: Retrieve Product By ID
    When I request to the Product that is created
    Then I should receive the details of the Product with status 200

  Scenario: Create a product with invalid data
    When I create a product with invalid data
    Then I should receive an error

  Scenario: Delete the product created
    When I request to delete the Product that is created
    Then the Product should be deleted with status 202