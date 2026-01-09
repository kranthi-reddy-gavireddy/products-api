Feature: Product Management
  As a user of the products API
  I want to manage products
  So that I can create and retrieve product information

  Scenario: Retrieve all products
    When I request all products
    Then I should receive a list of products

  Scenario: Create a new product
    When I create a product with name "Test Product", price 10.99, seller_id "seller123", quantity 5
    Then the product should be created successfully
    And I should receive the product details

  Scenario: Create a product with invalid data
    When I create a product with invalid data
    Then I should receive an error