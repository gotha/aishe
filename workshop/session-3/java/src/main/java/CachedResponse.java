import com.aishe.client.models.AnswerResponse;

/**
 * Wrapper for cached response with similarity score.
 */
public class CachedResponse {
    AnswerResponse response;
    Double similarity;

    CachedResponse(AnswerResponse response, Double similarity) {
        this.response = response;
        this.similarity = similarity;
    }
}

