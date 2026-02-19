import com.google.gson.annotations.SerializedName;
import java.util.List;

/**
 * Request/response models for LangCache REST API.
 */

class LangCacheSearchRequest {
    String prompt;
    @SerializedName("similarity_threshold")
    double similarityThreshold;
}

class LangCacheSearchEntry {
    String prompt;
    String response;
    Double similarity;
}

class LangCacheSearchResponse {
    List<LangCacheSearchEntry> data;
}

class LangCacheSetRequest {
    String prompt;
    String response;
}

