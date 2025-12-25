package com.sms.sender.service;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.List;

/**
 * Service for managing and checking the SMS block list stored in Redis.
 * Handles block list initialization and phone number validation.
 */
@Service
@Slf4j
public class BlockListService {

    private final RedisTemplate<String, String> redisTemplate;

    @Value("${app.redis.blocklist-key}")
    private String blocklistKey;

    /**
     * Dummy blocked phone numbers for testing.
     */
    private static final List<String> DUMMY_BLOCKED_NUMBERS = Arrays.asList(
            "+1111111111",
            "+2222222222",
            "+3333333333",
            "+9999999999",
            "+5555555555"
    );

    public BlockListService(RedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    /**
     * Initialize Redis with dummy blocked users on application startup.
     * This method runs automatically after bean construction.
     */
    @PostConstruct
    public void initializeBlockList() {
        try {
            log.info("Initializing block list with dummy data...");
            
            // Check if block list already exists
            Long existingCount = redisTemplate.opsForSet().size(blocklistKey);
            
            if (existingCount != null && existingCount > 0) {
                log.info("Block list already exists with {} entries. Skipping initialization.", existingCount);
                return;
            }
            
            // Add dummy blocked numbers to Redis SET
            for (String phoneNumber : DUMMY_BLOCKED_NUMBERS) {
                redisTemplate.opsForSet().add(blocklistKey, phoneNumber);
            }
            
            log.info("Successfully initialized block list with {} dummy blocked numbers: {}", 
                    DUMMY_BLOCKED_NUMBERS.size(), DUMMY_BLOCKED_NUMBERS);
            
        } catch (Exception e) {
            log.error("Failed to initialize block list in Redis: {}", e.getMessage(), e);
            // Don't throw exception to prevent application startup failure
        }
    }

    /**
     * Check if a phone number is in the block list.
     * 
     * @param phoneNumber Phone number to check
     * @return true if blocked, false otherwise
     */
    public boolean isBlocked(String phoneNumber) {
        try {
            Boolean isMember = redisTemplate.opsForSet().isMember(blocklistKey, phoneNumber);
            boolean blocked = Boolean.TRUE.equals(isMember);
            
            if (blocked) {
                log.warn("Phone number {} is in the block list", phoneNumber);
            } else {
                log.debug("Phone number {} is not blocked", phoneNumber);
            }
            
            return blocked;
            
        } catch (Exception e) {
            log.error("Error checking block list for phone number {}: {}", phoneNumber, e.getMessage(), e);
            // Default to not blocked in case of Redis failure
            return false;
        }
    }

    /**
     * Add a phone number to the block list.
     * 
     * @param phoneNumber Phone number to block
     * @return true if successfully added, false otherwise
     */
    public boolean addToBlockList(String phoneNumber) {
        try {
            Long result = redisTemplate.opsForSet().add(blocklistKey, phoneNumber);
            boolean added = result != null && result > 0;
            
            if (added) {
                log.info("Added phone number {} to block list", phoneNumber);
            } else {
                log.info("Phone number {} was already in block list", phoneNumber);
            }
            
            return added;
            
        } catch (Exception e) {
            log.error("Error adding phone number {} to block list: {}", phoneNumber, e.getMessage(), e);
            return false;
        }
    }

    /**
     * Remove a phone number from the block list.
     * 
     * @param phoneNumber Phone number to unblock
     * @return true if successfully removed, false otherwise
     */
    public boolean removeFromBlockList(String phoneNumber) {
        try {
            Long result = redisTemplate.opsForSet().remove(blocklistKey, phoneNumber);
            boolean removed = result != null && result > 0;
            
            if (removed) {
                log.info("Removed phone number {} from block list", phoneNumber);
            } else {
                log.info("Phone number {} was not in block list", phoneNumber);
            }
            
            return removed;
            
        } catch (Exception e) {
            log.error("Error removing phone number {} from block list: {}", phoneNumber, e.getMessage(), e);
            return false;
        }
    }

    /**
     * Get the total count of blocked phone numbers.
     * 
     * @return Number of blocked phone numbers
     */
    public long getBlockListSize() {
        try {
            Long size = redisTemplate.opsForSet().size(blocklistKey);
            return size != null ? size : 0;
        } catch (Exception e) {
            log.error("Error getting block list size: {}", e.getMessage(), e);
            return 0;
        }
    }
}
