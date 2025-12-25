package com.sms.sender.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

/**
 * Request model for sending SMS messages.
 * Contains the phone number and message content.
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class SmsRequest {

    /**
     * Phone number in international format (e.g., +1234567890)
     * Must start with + followed by digits, 10-15 digits total
     */
    @NotBlank(message = "Phone number is required")
    @Pattern(regexp = "^\\+?[1-9]\\d{9,14}$", message = "Invalid phone number format. Must be 10-15 digits.")
    @JsonProperty("phoneNumber")
    private String phoneNumber;

    /**
     * SMS message content
     * Must be between 1 and 160 characters
     */
    @NotBlank(message = "Message is required")
    @Size(min = 1, max = 160, message = "Message must be between 1 and 160 characters")
    @JsonProperty("message")
    private String message;
}
