package com.sms.sender.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * Response model for SMS sending operations.
 * Contains the status, message, and timestamp of the operation.
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class SmsResponse {

    /**
     * Status of the SMS sending operation
     * Possible values: SUCCESS, FAILED, BLOCKED
     */
    @JsonProperty("status")
    private String status;

    /**
     * Descriptive message about the operation result
     */
    @JsonProperty("message")
    private String message;

    /**
     * Timestamp when the operation was processed
     */
    @JsonProperty("timestamp")
    private LocalDateTime timestamp;

    /**
     * Phone number that was processed
     */
    @JsonProperty("phoneNumber")
    private String phoneNumber;

    /**
     * Helper method to create a success response
     */
    public static SmsResponse success(String phoneNumber, String message) {
        return SmsResponse.builder()
                .status("SUCCESS")
                .message(message)
                .phoneNumber(phoneNumber)
                .timestamp(LocalDateTime.now())
                .build();
    }

    /**
     * Helper method to create a failed response
     */
    public static SmsResponse failed(String phoneNumber, String message) {
        return SmsResponse.builder()
                .status("FAILED")
                .message(message)
                .phoneNumber(phoneNumber)
                .timestamp(LocalDateTime.now())
                .build();
    }

    /**
     * Helper method to create a blocked response
     */
    public static SmsResponse blocked(String phoneNumber) {
        return SmsResponse.builder()
                .status("BLOCKED")
                .message("Phone number is in the block list")
                .phoneNumber(phoneNumber)
                .timestamp(LocalDateTime.now())
                .build();
    }
}
