Feature: Order Management

  Scenario: Create a new order
    Given I have a valid order request
    When I send the order request to the order service
    Then I should receive a confirmation of the order creation

  Scenario: Retrieve an existing order
    Given I have an existing order with ID "12345"
    When I request the order details for ID "12345"
    Then I should receive the order details with status "PENDING"

  Scenario: Update an existing order
    Given I have an existing order with ID "12345"
    When I update the order status to "PENDING"
    Then I should receive a confirmation that the order status has been updated

  Scenario: Delete an existing order
    Given I have an existing order with ID "12345"
    When I delete the order with ID "12345"
    Then I should receive a confirmation that the order has been deleted
    And the order should no longer exist in the system
