package com.aishe.client;

/**
 * Configuration constants for AISHE client.
 */
public class Config {
    /** AISHE API server host */
    public static final String SERVER_HOST = getEnv("SERVER_HOST", "localhost");
    
    /** AISHE API server port */
    public static final int SERVER_PORT = getEnvInt("SERVER_PORT", 8000);
    
    /** AISHE API URL */
    public static final String AISHE_API_URL = getEnv("AISHE_API_URL", 
        String.format("http://%s:%d", SERVER_HOST, SERVER_PORT));
    
    /** Request timeout in milliseconds */
    public static final int REQUEST_TIMEOUT_MS = getEnvInt("REQUEST_TIMEOUT", 120) * 1000;
    
    /** Redis host */
    public static final String REDIS_HOST = getEnv("REDIS_HOST", "localhost");
    
    /** Redis port */
    public static final int REDIS_PORT = getEnvInt("REDIS_PORT", 6379);
    
    /** Redis database */
    public static final int REDIS_DATABASE = getEnvInt("REDIS_DATABASE", 0);
    
    /** Redis username */
    public static final String REDIS_USERNAME = getEnv("REDIS_USERNAME", "default");
    
    /** Redis password */
    public static final String REDIS_PASSWORD = getEnv("REDIS_PASSWORD", "");
    
    /** Redis cache key prefix */
    public static final String REDIS_CACHE_KEY_PREFIX = getEnv("REDIS_CACHE_KEY_PREFIX", "aishe:question");
    
    /**
     * Get environment variable as string with default value.
     */
    private static String getEnv(String key, String defaultValue) {
        String value = System.getenv(key);
        return value != null ? value : defaultValue;
    }
    
    /**
     * Get environment variable as integer with default value.
     */
    private static int getEnvInt(String key, int defaultValue) {
        String value = System.getenv(key);
        if (value == null) {
            return defaultValue;
        }
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            return defaultValue;
        }
    }
    
    /**
     * Display the current configuration.
     */
    public static void displayConfig() {
        System.out.println("Configuration:");
        System.out.println("  Server Host: " + SERVER_HOST);
        System.out.println("  Server Port: " + SERVER_PORT);
        System.out.println("  API URL: " + AISHE_API_URL);
        System.out.println("  Request Timeout: " + REQUEST_TIMEOUT_MS + "ms");
        System.out.println("  Redis Host: " + REDIS_HOST);
        System.out.println("  Redis Port: " + REDIS_PORT);
        System.out.println("  Redis Database: " + REDIS_DATABASE);
        System.out.println("  Redis Username: " + REDIS_USERNAME);
        System.out.println("  Redis Password: " + REDIS_PASSWORD);
        System.out.println("  Redis Cache Key Prefix: " + REDIS_CACHE_KEY_PREFIX);
    }
}

