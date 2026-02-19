package com.aishe.client;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * Utility functions for AISHE client.
 */
public class Utils {
    /**
     * Generate a cache key for a question.
     * 
     * Uses SHA-256 hash to generate a unique key.
     *
     * @param question Question to generate a cache key for
     * @return Cache key in format: {prefix}:{hash}
     */
    public static String generateCacheKey(String question) {
        return generateCacheKey(question, Config.REDIS_CACHE_KEY_PREFIX);
    }

    /**
     * Generate a cache key for a question with custom prefix.
     * 
     * Uses SHA-256 hash to generate a unique key.
     *
     * @param question Question to generate a cache key for
     * @param prefix Cache key prefix
     * @return Cache key in format: {prefix}:{hash}
     */
    public static String generateCacheKey(String question, String prefix) {
        try {
            // Normalize the question (lowercase, trim whitespace)
            String normalized = question.toLowerCase().trim();
            
            // Create SHA-256 hash
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(normalized.getBytes(StandardCharsets.UTF_8));
            
            // Convert to hex string
            StringBuilder hexString = new StringBuilder();
            for (byte b : hash) {
                String hex = Integer.toHexString(0xff & b);
                if (hex.length() == 1) {
                    hexString.append('0');
                }
                hexString.append(hex);
            }
            
            return prefix + ":" + hexString.toString();
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("SHA-256 algorithm not available", e);
        }
    }
}

