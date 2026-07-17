// SlidingWindow allows at most `limit` requests in any window_ms
// span, measured from each request.
//
// TODO: this is the fixed-window counter this exercise exists to
// replace -- it resets at boundaries, so a burst on each side of one
// puts 2x the limit through.
class SlidingWindow {
public:
    SlidingWindow(int limit, long window_ms)
        : limit_(limit), window_ms_(window_ms) {}

    bool AllowAt(long now_ms) {
        if (now_ms - window_start_ >= window_ms_) {
            window_start_ = now_ms;
            count_ = 0;
        }
        if (count_ < limit_) {
            count_++;
            return true;
        }
        return false;
    }

private:
    int limit_;
    long window_ms_;
    long window_start_ = 0;
    int count_ = 0;
};
