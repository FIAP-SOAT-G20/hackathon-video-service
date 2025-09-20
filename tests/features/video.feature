Feature: Video Management

  Scenario: Create a new video
    Given I have a valid video request
    When I send the video request to the video service
    Then I should receive a confirmation of the video creation

  Scenario: Retrieve an existing video
    Given I have an existing video with ID "12345"
    When I request the video details for ID "12345"
    Then I should receive the video details with status "PROCESSING"

  Scenario: Update an existing video
    Given I have an existing video with ID "12345"
    When I update the video status to "FINISHED"
    Then I should receive a confirmation that the video status has been updated

  Scenario: Delete an existing video
    Given I have an existing video with ID "12345"
    When I delete the video with ID "12345"
    Then I should receive a confirmation that the video has been deleted
    And the video should no longer exist in the system
