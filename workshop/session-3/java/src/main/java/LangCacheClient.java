/**
 * Simple client holder for LangCache API credentials.
 */
public class LangCacheClient {
    String serverUrl;
    String cacheId;
    String apiKey;

    LangCacheClient(String serverUrl, String cacheId, String apiKey) {
        this.serverUrl = serverUrl;
        this.cacheId = cacheId;
        this.apiKey = apiKey;
    }
}

