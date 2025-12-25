package com.sms.sender.kafka;

import com.sms.sender.model.KafkaEvent;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;

import java.util.concurrent.CompletableFuture;

/**
 * Service for producing SMS events to Kafka topic.
 * NOTE: For v0, only synchronous message publishing (sendSmsEventSync) is used
 * to ensure message delivery confirmation before responding to clients.
 * Asynchronous method is commented out for future use.
 */
@Service
@Slf4j
public class SmsKafkaProducer {

    private final KafkaTemplate<String, Object> kafkaTemplate;

    @Value("${app.kafka.topic}")
    private String kafkaTopic;

    public SmsKafkaProducer(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    /*
     * ASYNC METHOD - COMMENTED OUT FOR v0
     * 
     * For v0, we use synchronous sending only (sendSmsEventSync) to ensure
     * Kafka message delivery before responding to the client.
     * This async method can be enabled in future versions for better performance.
     * 
    /**
     * Send an SMS event to Kafka asynchronously.
     * Uses phone number as the message key for partitioning.
     * 
     * @param event The KafkaEvent to send
     *
    public void sendSmsEvent(KafkaEvent event) {
        try {
            log.debug("Sending SMS event to Kafka topic '{}': eventId={}, phoneNumber={}, status={}", 
                    kafkaTopic, event.getEventId(), event.getPhoneNumber(), event.getStatus());
            
            // Use phone number as key for consistent partitioning
            String key = event.getPhoneNumber();
            
            // Send message asynchronously
            CompletableFuture<SendResult<String, Object>> future = 
                    kafkaTemplate.send(kafkaTopic, key, event);
            
            // Add callback for success/failure handling
            future.whenComplete((result, ex) -> {
                if (ex == null) {
                    log.info("Successfully sent SMS event to Kafka: eventId={}, phoneNumber={}, partition={}, offset={}", 
                            event.getEventId(), 
                            event.getPhoneNumber(),
                            result.getRecordMetadata().partition(),
                            result.getRecordMetadata().offset());
                } else {
                    log.error("Failed to send SMS event to Kafka: eventId={}, phoneNumber={}, error={}", 
                            event.getEventId(), 
                            event.getPhoneNumber(), 
                            ex.getMessage(), 
                            ex);
                }
            });
            
        } catch (Exception e) {
            log.error("Exception while sending SMS event to Kafka: eventId={}, phoneNumber={}, error={}", 
                    event.getEventId(), 
                    event.getPhoneNumber(), 
                    e.getMessage(), 
                    e);
        }
    }
    */

    /**
     * Send an SMS event to Kafka synchronously and wait for confirmation.
     * 
     * NOTE: This is the PRIMARY method used in v0 to ensure reliable message delivery
     * before responding to the client. Provides confirmation that the event was successfully
     * written to Kafka before the API returns.
     * 
     * @param event The KafkaEvent to send
     * @return true if successfully sent, false otherwise
     */
    public boolean sendSmsEventSync(KafkaEvent event) {
        try {
            log.debug("Sending SMS event to Kafka (sync) topic '{}': eventId={}, phoneNumber={}", 
                    kafkaTopic, event.getEventId(), event.getPhoneNumber());
            
            String key = event.getPhoneNumber();
            SendResult<String, Object> result = kafkaTemplate.send(kafkaTopic, key, event).get();
            
            log.info("Successfully sent SMS event to Kafka (sync): eventId={}, partition={}, offset={}", 
                    event.getEventId(),
                    result.getRecordMetadata().partition(),
                    result.getRecordMetadata().offset());
            
            return true;
            
        } catch (Exception e) {
            log.error("Failed to send SMS event to Kafka (sync): eventId={}, error={}", 
                    event.getEventId(), 
                    e.getMessage(), 
                    e);
            return false;
        }
    }
}
